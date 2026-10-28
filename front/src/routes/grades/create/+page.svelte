<script lang="ts">
    import { onMount } from 'svelte';
    import { resolve } from '$app/paths';
    import { PUBLIC_API_URL } from '$env/static/public';

    let testId = $state<string | number>('');
    let enrollmentId = $state<string | number>('');
    let value = $state<string | number>('');

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
            const response = await fetch(`${PUBLIC_API_URL}/grades`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    test_id: Number(testId),
                    enrollment_id: Number(enrollmentId),
                    value: Number(value)
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to create grade');
            }

            window.location.href = '/grades';
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:head>
    <title>New Grade | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href={resolve('/grades')} class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Grades
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Create Grade</h1>

            <form onsubmit={handleCreate} class="flex flex-col gap-5">
                <div>
                    <label for="test-id" class="block text-sm font-medium text-gray-700 mb-1">Test ID</label>
                    <input id="test-id" type="number" bind:value={testId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="enrollment-id" class="block text-sm font-medium text-gray-700 mb-1">Enrollment ID</label>
                    <input id="enrollment-id" type="number" bind:value={enrollmentId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="value" class="block text-sm font-medium text-gray-700 mb-1">Grade Value (0-10)</label>
                    <input id="value" type="number" step="0.1" bind:value={value} required min="0" max="10" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                {#if errorMessage}
                    <div class="bg-red-50 text-red-500 text-sm p-3 rounded-lg text-center">{errorMessage}</div>
                {/if}

                <button type="submit" disabled={isSaving} class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg mt-2 transition-colors">
                    {isSaving ? 'Creating...' : 'Create Grade'}
                </button>
            </form>
        </div>
    </div>
</div>
