<script lang="ts">
    import { onMount } from 'svelte';
    import { resolve } from '$app/paths';

    let studentId = $state<string | number>('');
    let courseId = $state<string | number>('');
    let average = $state<number | null>(null);

    let isLoading = $state(false);
    let errorMessage = $state('');

    onMount(() => {
        if (!localStorage.getItem('paseto_token')) {
            window.location.href = '/login';
        }
    });

    async function handleCalculate(event: SubmitEvent) {
        event.preventDefault();
        isLoading = true;
        errorMessage = '';
        average = null;
        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch(`http://localhost:8080/grades/average/${studentId}/${courseId}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to calculate average');
            }

            const data = await response.json();
            average = data.data ? data.data.average : data.average;
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isLoading = false;
        }
    }
</script>

<svelte:head>
    <title>Calculate Average | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href={resolve('/grades/averages')} class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Report Cards
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Calculate Average</h1>

            <form onsubmit={handleCalculate} class="flex flex-col gap-5">
                <div>
                    <label for="student-id" class="block text-sm font-medium text-gray-700 mb-1">Student ID</label>
                    <input id="student-id" type="number" bind:value={studentId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="course-id" class="block text-sm font-medium text-gray-700 mb-1">Course ID</label>
                    <input id="course-id" type="number" bind:value={courseId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                {#if errorMessage}
                    <div class="bg-red-50 text-red-500 text-sm p-3 rounded-lg text-center">{errorMessage}</div>
                {/if}

                {#if average !== null}
                    <div class="bg-green-50 text-green-700 text-lg font-bold p-4 rounded-lg text-center mt-2">
                        Average: {average.toFixed(2)}
                    </div>
                {/if}

                <button type="submit" disabled={isLoading} class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg mt-2 transition-colors">
                    {isLoading ? 'Calculating...' : 'Calculate'}
                </button>
            </form>
        </div>
    </div>
</div>
