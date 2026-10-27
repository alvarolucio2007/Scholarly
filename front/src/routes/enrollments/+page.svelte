<script lang="ts">
    import { onMount } from 'svelte';
    import { resolve } from '$app/paths';

    let enrollments = $state<any[]>([]);
    let isLoading = $state(true);
    let errorMessage = $state('');

    async function fetchEnrollments() {
        isLoading = true;
        errorMessage = '';
        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch('http://localhost:8080/enrollments', {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (response.status === 401) {
                localStorage.removeItem('paseto_token');
                window.location.href = '/login';
                return;
            }

            if (!response.ok) throw new Error('Failed to fetch enrollments');

            const data = await response.json();
            enrollments = data.data ? data.data : data;

        } catch (error) {
            errorMessage = 'Could not load enrollments.';
        } finally {
            isLoading = false;
        }
    }

    async function deleteEnrollment(id: number) {
        if (!confirm('Are you sure you want to delete this enrollment?')) return;

        const token = localStorage.getItem('paseto_token');
        try {
            const response = await fetch(`http://localhost:8080/enrollments/${id}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (response.ok) {
                enrollments = enrollments.filter(e => e.id !== id);
            } else {
                alert('Failed to delete enrollment.');
            }
        } catch (error) {
            alert('Connection error.');
        }
    }

    onMount(() => {
        if (!localStorage.getItem('paseto_token')) {
            window.location.href = '/login';
            return;
        }
        fetchEnrollments();
    });
</script>

<svelte:head>
    <title>Enrollments | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8">
    <header class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
            <a href={resolve('/dashboard')} class="text-blue-600 hover:underline text-sm font-medium mb-2 inline-block">
                &larr; Back to Dashboard
            </a>
            <h1 class="text-3xl font-bold text-gray-800">Enrollments Management</h1>
            <p class="text-gray-500">View and manage student enrollments</p>
        </div>

        <a href={resolve('/courses/enroll')} class="bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-6 rounded-lg transition-colors inline-block">
            + New Enrollment
        </a>
    </header>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {#if isLoading}
            <div class="p-10 text-center text-gray-500 font-medium">Loading enrollments...</div>
        {:else if errorMessage}
            <div class="p-10 text-center text-red-500 font-medium">{errorMessage}</div>
        {:else if enrollments.length === 0}
            <div class="p-10 text-center text-gray-500 font-medium">No enrollments found.</div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-gray-50 border-b border-gray-200 text-gray-600 text-sm">
                            <th class="p-4 font-semibold">ID</th>
                            <th class="p-4 font-semibold">Student ID</th>
                            <th class="p-4 font-semibold">Course ID</th>
                            <th class="p-4 font-semibold">Date</th>
                            <th class="p-4 font-semibold">Status</th>
                            <th class="p-4 font-semibold text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100">
                        {#each enrollments as enrollment}
                            <tr class="hover:bg-gray-50 transition-colors">
                                <td class="p-4 text-gray-500 font-mono text-sm">#{enrollment.id}</td>
                                <td class="p-4 text-gray-600">#{enrollment.student_id}</td>
                                <td class="p-4 text-gray-600">#{enrollment.course_id}</td>
                                <td class="p-4 text-gray-600">{new Date(enrollment.enrolled_at).toLocaleDateString()}</td>
                                <td class="p-4">
                                    <span class="px-2 py-1 rounded text-xs font-semibold
                                        {enrollment.status === 'active' ? 'bg-green-100 text-green-700' : ''}
                                        {enrollment.status === 'cancelled' ? 'bg-red-100 text-red-700' : ''}
                                        {enrollment.status === 'completed' ? 'bg-blue-100 text-blue-700' : ''}">
                                        {enrollment.status.toUpperCase()}
                                    </span>
                                </td>
                                <td class="p-4 text-right flex justify-end gap-2">
                                    <a
                                        href={resolve('/enrollments/[id]', { id: enrollment.id.toString() })}
                                        class="text-blue-600 hover:bg-blue-50 px-3 py-1 rounded-md text-sm font-medium transition-colors"
                                    >
                                        Edit
                                    </a>
                                    <button
                                        onclick={() => deleteEnrollment(enrollment.id)}
                                        class="text-red-600 hover:bg-red-50 px-3 py-1 rounded-md text-sm font-medium transition-colors"
                                    >
                                        Delete
                                    </button>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
</div>
