<script lang="ts">
    import { onMount } from 'svelte';

    let userId = $state<string | number>('');
    let department = $state('');
    let isSaving = $state(false);
    let errorMessage = $state('');

    onMount(() => {
        if (!localStorage.getItem('paseto_token')) {
            window.location.href = '/login';
        }
    });

    async function handleCreate(event: SubmitEvent) {
        event.preventDefault();
        isSaving = true;
        errorMessage = '';
        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch('http://localhost:8080/teachers', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    user_id: Number(userId),
                    department: department
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to create teacher');
            }

            window.location.href = '/teachers';
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:head>
    <title>New Teacher | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href="/teachers" class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Teachers
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Create Teacher</h1>

            <form onsubmit={handleCreate} class="flex flex-col gap-5">
                <div>
                    <label for="user-id" class="block text-sm font-medium text-gray-700 mb-1">User ID</label>
                    <input
                        id="user-id"
                        type="number"
                        bind:value={userId}
                        required
                        min="1"
                        class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                    />
                </div>

                <div>
                    <label for="department" class="block text-sm font-medium text-gray-700 mb-1">Department</label>
                    <input
                        id="department"
                        type="text"
                        bind:value={department}
                        required
                        maxlength="100"
                        class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none"
                    />
                </div>

                {#if errorMessage}
                    <div class="bg-red-50 text-red-500 text-sm p-3 rounded-lg text-center">{errorMessage}</div>
                {/if}

                <button
                    type="submit"
                    disabled={isSaving}
                    class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg mt-2 transition-colors"
                >
                    {isSaving ? 'Creating...' : 'Create Teacher'}
                </button>
            </form>
        </div>
    </div>
</div>
