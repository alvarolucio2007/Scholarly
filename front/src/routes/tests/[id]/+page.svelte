<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/state';
    import { resolve } from '$app/paths';
    import { PUBLIC_API_URL } from '$env/static/public';

    let testId = $derived(page.params.id);

    let courseId = $state<string | number>('');
    let name = $state('');
    let weight = $state<string | number>('');
    let testDate = $state('');

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
            const response = await fetch(`${PUBLIC_API_URL}/tests/${testId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (!response.ok) throw new Error('Test not found');

            const data = await response.json();
            const test = data.data ? data.data : data;

            courseId = test.course_id;
            name = test.name;
            weight = test.weight;

            if (test.test_date) {
                testDate = new Date(test.test_date).toISOString().slice(0, 16);
            }
        } catch (err) {
            errorMessage = 'Failed to load test data.';
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
            const isoDate = testDate ? new Date(testDate).toISOString() : '';

            const response = await fetch(`${PUBLIC_API_URL}/tests`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    id: Number(testId),
                    course_id: Number(courseId),
                    name: name.trim(),
                    weight: Number(weight),
                    test_date: isoDate
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to update test');
            }

            successMessage = 'Test updated successfully!';
            setTimeout(() => { window.location.href = '/tests'; }, 1500);
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:head>
    <title>Edit Test | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href={resolve('/tests')} class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Tests
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Edit Test #{testId}</h1>

            {#if isLoading}
                <p class="text-gray-500">Loading test data...</p>
            {:else}
                <form onsubmit={handleUpdate} class="flex flex-col gap-5">
                    <div>
                        <label for="edit-course-id" class="block text-sm font-medium text-gray-700 mb-1">Course ID</label>
                        <input id="edit-course-id" type="number" bind:value={courseId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
                        <input id="edit-name" type="text" bind:value={name} required maxlength="255" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-weight" class="block text-sm font-medium text-gray-700 mb-1">Weight</label>
                        <input id="edit-weight" type="number" step="0.1" bind:value={weight} required min="0.1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-test-date" class="block text-sm font-medium text-gray-700 mb-1">Test Date & Time</label>
                        <input id="edit-test-date" type="datetime-local" bind:value={testDate} required class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
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
