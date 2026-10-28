<script lang="ts">
    import { onMount } from 'svelte';
    import { resolve } from '$app/paths';
    import { PUBLIC_API_URL } from '$env/static/public';

    let courseId = $state<string | number>('');
    let name = $state('');
    let weight = $state<string | number>('');
    let testDate = $state('');

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
            const isoDate = new Date(testDate).toISOString();

            const response = await fetch(`${PUBLIC_API_URL}/tests`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    course_id: Number(courseId),
                    name: name.trim(),
                    weight: Number(weight),
                    test_date: isoDate
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to create test');
            }

            window.location.href = '/tests';
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:head>
    <title>New Test | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href={resolve('/tests')} class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Tests
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Create Test</h1>

            <form onsubmit={handleCreate} class="flex flex-col gap-5">
                <div>
                    <label for="course-id" class="block text-sm font-medium text-gray-700 mb-1">Course ID</label>
                    <input id="course-id" type="number" bind:value={courseId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
                    <input id="name" type="text" bind:value={name} required maxlength="255" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="weight" class="block text-sm font-medium text-gray-700 mb-1">Weight</label>
                    <input id="weight" type="number" step="0.1" bind:value={weight} required min="0.1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="test-date" class="block text-sm font-medium text-gray-700 mb-1">Test Date & Time</label>
                    <input id="test-date" type="datetime-local" bind:value={testDate} required class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                {#if errorMessage}
                    <div class="bg-red-50 text-red-500 text-sm p-3 rounded-lg text-center">{errorMessage}</div>
                {/if}

                <button type="submit" disabled={isSaving} class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg mt-2 transition-colors">
                    {isSaving ? 'Creating...' : 'Create Test'}
                </button>
            </form>
        </div>
    </div>
</div>
