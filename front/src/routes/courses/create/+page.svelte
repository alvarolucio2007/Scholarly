<script lang="ts">
    import { onMount } from 'svelte';
    import { PUBLIC_API_URL } from '$env/static/public';

    let teacherId = $state<string | number>('');
    let name = $state('');
    let code = $state('');
    let semester = $state('');
    let maxStudents = $state<string | number>('');

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
            const response = await fetch(`${PUBLIC_API_URL}/courses`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    teacher_id: Number(teacherId),
                    name,
                    code,
                    semester,
                    max_students: Number(maxStudents)
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to create course');
            }

            window.location.href = '/courses';
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:head>
    <title>New Course | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href="/courses" class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Courses
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Create Course</h1>

            <form onsubmit={handleCreate} class="flex flex-col gap-5">
                <div>
                    <label for="teacher-id" class="block text-sm font-medium text-gray-700 mb-1">Teacher ID</label>
                    <input id="teacher-id" type="number" bind:value={teacherId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
                    <input id="name" type="text" bind:value={name} required maxlength="255" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="code" class="block text-sm font-medium text-gray-700 mb-1">Code</label>
                    <input id="code" type="text" bind:value={code} required maxlength="50" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="semester" class="block text-sm font-medium text-gray-700 mb-1">Semester</label>
                    <input id="semester" type="text" bind:value={semester} required maxlength="50" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                <div>
                    <label for="max-students" class="block text-sm font-medium text-gray-700 mb-1">Max Students</label>
                    <input id="max-students" type="number" bind:value={maxStudents} required min="1" max="32767" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                </div>

                {#if errorMessage}
                    <div class="bg-red-50 text-red-500 text-sm p-3 rounded-lg text-center">{errorMessage}</div>
                {/if}

                <button type="submit" disabled={isSaving} class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg mt-2 transition-colors">
                    {isSaving ? 'Creating...' : 'Create Course'}
                </button>
            </form>
        </div>
    </div>
</div>
