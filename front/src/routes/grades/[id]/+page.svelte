<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/state';
    import { resolve } from '$app/paths';
    import { PUBLIC_API_URL } from '$env/static/public';

    let gradeId = $derived(page.params.id);

    let testId = $state<string | number>('');
    let enrollmentId = $state<string | number>('');
    let value = $state<string | number>('');

    let isLoading = $state(true);
    let isSaving = $state(false);
    let errorMessage = $state('');
    let successMessage = $state('');

    onMount(async () => {
        const token = localStorage.getItem('paseto_token');
        if (!token) {
            window.location.href = '/login';
            return;
        }

        try {
            const response = await fetch(`${PUBLIC_API_URL}/grades/${gradeId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (!response.ok) throw new Error('Grade not found');

            const data = await response.json();
            const grade = data.data ? data.data : data;

            testId = grade.test_id;
            enrollmentId = grade.enrollment_id;
            value = grade.value;
        } catch (err) {
            errorMessage = 'Failed to load grade data.';
        } finally {
            isLoading = false;
        }
    });

    async function handleUpdate(event: SubmitEvent) {
        event.preventDefault();
        isSaving = true;
        errorMessage = '';
        successMessage = '';

        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch(`${PUBLIC_API_URL}/grades`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    id: Number(gradeId),
                    test_id: Number(testId),
                    enrollment_id: Number(enrollmentId),
                    value: Number(value)
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to update grade');
            }

            successMessage = 'Grade updated successfully!';
            setTimeout(() => { window.location.href = '/grades'; }, 1500);
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:head>
    <title>Edit Grade | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href={resolve('/grades')} class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Grades
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Edit Grade #{gradeId}</h1>

            {#if isLoading}
                <p class="text-gray-500">Loading grade data...</p>
            {:else}
                <form onsubmit={handleUpdate} class="flex flex-col gap-5">
                    <div>
                        <label for="edit-test-id" class="block text-sm font-medium text-gray-700 mb-1">Test ID</label>
                        <input id="edit-test-id" type="number" bind:value={testId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-enrollment-id" class="block text-sm font-medium text-gray-700 mb-1">Enrollment ID</label>
                        <input id="edit-enrollment-id" type="number" bind:value={enrollmentId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-value" class="block text-sm font-medium text-gray-700 mb-1">Grade Value (0-10)</label>
                        <input id="edit-value" type="number" step="0.1" bind:value={value} required min="0" max="10" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    {#if errorMessage}
                        <div class="bg-red-50 text-red-500 text-sm p-3 rounded-lg text-center">{errorMessage}</div>
                    {/if}

                    {#if successMessage}
                        <div class="bg-green-50 text-green-600 text-sm p-3 rounded-lg text-center">{successMessage}</div>
                    {/if}

                    <button type="submit" disabled={isSaving} class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg mt-2 transition-colors">
                        {isSaving ? 'Saving...' : 'Save Changes'}
                    </button>
                </form>
            {/if}
        </div>
    </div>
</div>
