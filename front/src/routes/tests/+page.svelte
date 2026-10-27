<script lang="ts">
    import { onMount } from 'svelte';
    import { resolve } from '$app/paths';

    let tests = $state<any[]>([]);
    let isLoading = $state(true);
    let errorMessage = $state('');

    async function fetchTests() {
        isLoading = true;
        errorMessage = '';
        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch('http://localhost:8080/tests', {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (response.status === 401) {
                localStorage.removeItem('paseto_token');
                window.location.href = '/login';
                return;
            }

            if (!response.ok) throw new Error('Failed to fetch tests');

            const data = await response.json();
            tests = data.data ? data.data : data;

        } catch (error) {
            errorMessage = 'Could not load tests.';
        } finally {
            isLoading = false;
        }
    }

    async function deleteTest(id: number) {
        if (!confirm('Are you sure you want to delete this test?')) return;

        const token = localStorage.getItem('paseto_token');
        try {
            const response = await fetch(`http://localhost:8080/tests/${id}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (response.ok) {
                tests = tests.filter(t => t.id !== id);
            } else {
                alert('Failed to delete test.');
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
        fetchTests();
    });
</script>

<svelte:head>
    <title>Tests | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8">
    <header class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
            <a href={resolve('/dashboard')} class="text-blue-600 hover:underline text-sm font-medium mb-2 inline-block">
                &larr; Back to Dashboard
            </a>
            <h1 class="text-3xl font-bold text-gray-800">Tests Management</h1>
            <p class="text-gray-500">View and manage course tests</p>
        </div>

        <a href={resolve('/tests/create')} class="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-6 rounded-lg transition-colors inline-block">
            + New Test
        </a>
    </header>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {#if isLoading}
            <div class="p-10 text-center text-gray-500 font-medium">Loading tests...</div>
        {:else if errorMessage}
            <div class="p-10 text-center text-red-500 font-medium">{errorMessage}</div>
        {:else if tests.length === 0}
            <div class="p-10 text-center text-gray-500 font-medium">No tests found.</div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-gray-50 border-b border-gray-200 text-gray-600 text-sm">
                            <th class="p-4 font-semibold">ID</th>
                            <th class="p-4 font-semibold">Course ID</th>
                            <th class="p-4 font-semibold">Name</th>
                            <th class="p-4 font-semibold">Weight</th>
                            <th class="p-4 font-semibold">Date</th>
                            <th class="p-4 font-semibold text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100">
                        {#each tests as test}
                            <tr class="hover:bg-gray-50 transition-colors">
                                <td class="p-4 text-gray-500 font-mono text-sm">#{test.id}</td>
                                <td class="p-4 text-gray-600">#{test.course_id}</td>
                                <td class="p-4 text-gray-800 font-medium">{test.name}</td>
                                <td class="p-4 text-gray-600">{test.weight}</td>
                                <td class="p-4 text-gray-600">{new Date(test.test_date).toLocaleDateString()}</td>
                                <td class="p-4 text-right flex justify-end gap-2">
                                    <a
                                        href={resolve('/tests/[id]', { id: test.id.toString() })}
                                        class="text-blue-600 hover:bg-blue-50 px-3 py-1 rounded-md text-sm font-medium transition-colors"
                                    >
                                        Edit
                                    </a>
                                    <button
                                        onclick={() => deleteTest(test.id)}
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
