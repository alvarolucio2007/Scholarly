<script lang="ts">
    import { onMount } from 'svelte';

    let teachers = $state<any[]>([]);
    let isLoading = $state(true);
    let errorMessage = $state('');

    async function fetchTeachers() {
        isLoading = true;
        errorMessage = '';
        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch('http://localhost:8080/teachers', {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (response.status === 401) {
                localStorage.removeItem('paseto_token');
                window.location.href = '/login';
                return;
            }

            if (!response.ok) throw new Error('Failed to fetch teachers');

            const data = await response.json();
            teachers = data.data ? data.data : data;

        } catch (error) {
            errorMessage = 'Could not load teachers.';
        } finally {
            isLoading = false;
        }
    }

    async function deleteTeacher(userId: number) {
        if (!confirm('Are you sure you want to delete this teacher?')) return;

        const token = localStorage.getItem('paseto_token');
        try {
            const response = await fetch(`http://localhost:8080/teachers/${userId}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (response.ok) {
                teachers = teachers.filter(t => t.user_id !== userId);
            } else {
                alert('Failed to delete teacher.');
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
        fetchTeachers();
    });
</script>

<svelte:head>
    <title>Teachers | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8">
    <header class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
            <a href="/dashboard" class="text-blue-600 hover:underline text-sm font-medium mb-2 inline-block">
                &larr; Back to Dashboard
            </a>
            <h1 class="text-3xl font-bold text-gray-800">Teachers Management</h1>
            <p class="text-gray-500">View and manage teaching staff</p>
        </div>

        <a href="/teachers/create" class="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-6 rounded-lg transition-colors inline-block">
            + New Teacher
        </a>
    </header>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {#if isLoading}
            <div class="p-10 text-center text-gray-500 font-medium">Loading teachers...</div>
        {:else if errorMessage}
            <div class="p-10 text-center text-red-500 font-medium">{errorMessage}</div>
        {:else if teachers.length === 0}
            <div class="p-10 text-center text-gray-500 font-medium">No teachers found.</div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-gray-50 border-b border-gray-200 text-gray-600 text-sm">
                            <th class="p-4 font-semibold">User ID</th>
                            <th class="p-4 font-semibold">Department</th>
                            <th class="p-4 font-semibold text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100">
                        {#each teachers as teacher}
                            <tr class="hover:bg-gray-50 transition-colors">
                                <td class="p-4 text-gray-500 font-mono text-sm">#{teacher.user_id}</td>
                                <td class="p-4 text-gray-800 font-medium">{teacher.department}</td>
                                <td class="p-4 text-right flex justify-end gap-2">
                                    <a
                                        href={`/teachers/${teacher.user_id}`}
                                        class="text-blue-600 hover:bg-blue-50 px-3 py-1 rounded-md text-sm font-medium transition-colors"
                                    >
                                        Edit
                                    </a>
                                    <button
                                        onclick={() => deleteTeacher(teacher.user_id)}
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
