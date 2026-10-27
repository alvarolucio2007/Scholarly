<script lang="ts">
    import { onMount } from 'svelte';
    import { resolve } from '$app/paths';

    let averages = $state<any[]>([]);
    let isLoading = $state(true);
    let errorMessage = $state('');

    async function fetchAverages() {
        isLoading = true;
        errorMessage = '';
        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch('http://localhost:8080/grades/average/list', {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (response.status === 401) {
                localStorage.removeItem('paseto_token');
                window.location.href = '/login';
                return;
            }

            if (!response.ok) throw new Error('Failed to fetch averages');

            const data = await response.json();
            averages = data.data ? data.data : data;

        } catch (error) {
            errorMessage = 'Could not load averages.';
        } finally {
            isLoading = false;
        }
    }

    onMount(() => {
        if (!localStorage.getItem('paseto_token')) {
            window.location.href = '/login';
            return;
        }
        fetchAverages();
    });
</script>

<svelte:head>
    <title>Report Cards | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8">
    <header class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
            <a href={resolve('/grades')} class="text-blue-600 hover:underline text-sm font-medium mb-2 inline-block">
                &larr; Back to Grades
            </a>
            <h1 class="text-3xl font-bold text-gray-800">Report Cards</h1>
            <p class="text-gray-500">Consolidated averages for students</p>
        </div>

        <a href={resolve('/grades/averages/calculate')} class="border-2 border-purple-600 text-purple-600 hover:bg-purple-50 font-semibold py-2 px-6 rounded-lg transition-colors inline-block">
            Calculate Specific Average
        </a>
    </header>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {#if isLoading}
            <div class="p-10 text-center text-gray-500 font-medium">Loading report cards...</div>
        {:else if errorMessage}
            <div class="p-10 text-center text-red-500 font-medium">{errorMessage}</div>
        {:else if averages.length === 0}
            <div class="p-10 text-center text-gray-500 font-medium">No report cards found.</div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-gray-50 border-b border-gray-200 text-gray-600 text-sm">
                            <th class="p-4 font-semibold">Student</th>
                            <th class="p-4 font-semibold">Course</th>
                            <th class="p-4 font-semibold">Teacher</th>
                            <th class="p-4 font-semibold">Average</th>
                            <th class="p-4 font-semibold">Status</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100">
                        {#each averages as card}
                            <tr class="hover:bg-gray-50 transition-colors">
                                <td class="p-4 text-gray-800 font-medium">{card.student_name} <span class="text-gray-400 text-xs ml-1">#{card.student_id}</span></td>
                                <td class="p-4 text-gray-600">{card.course_name} <span class="text-gray-400 text-xs">({card.course_code})</span></td>
                                <td class="p-4 text-gray-600">{card.teacher_name}</td>
                                <td class="p-4 text-gray-800 font-bold">{card.average.toFixed(1)}</td>
                                <td class="p-4">
                                    <span class="px-2 py-1 rounded text-xs font-semibold
                                        {card.situation === 'approved' ? 'bg-green-100 text-green-700' : ''}
                                        {card.situation === 'recovery' ? 'bg-yellow-100 text-yellow-700' : ''}
                                        {card.situation === 'cancelled' ? 'bg-red-100 text-red-700' : ''}">
                                        {card.situation.toUpperCase()}
                                    </span>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
</div>
