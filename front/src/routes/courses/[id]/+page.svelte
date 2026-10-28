<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/state';
    import { PUBLIC_API_URL } from '$env/static/public';

    let courseId = $derived(page.params.id);

    let teacherId = $state<string | number>('');
    let name = $state('');
    let code = $state('');
    let semester = $state('');
    let maxStudents = $state<string | number>('');

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
            const response = await fetch(`${PUBLIC_API_URL}/courses/${courseId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (!response.ok) throw new Error('Course not found');

            const data = await response.json();
            const course = data.data ? data.data : data;

            teacherId = course.teacher_id;
            name = course.name;
            code = course.code;
            semester = course.semester;
            maxStudents = course.max_students;
        } catch (err) {
            errorMessage = 'Failed to load course data.';
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
            const response = await fetch(`${PUBLIC_API_URL}/courses`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    id: Number(courseId),
                    teacher_id: Number(teacherId),
                    name,
                    code,
                    semester,
                    max_students: Number(maxStudents)
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to update course');
            }

            successMessage = 'Course updated successfully!';
            setTimeout(() => { window.location.href = '/courses'; }, 1500);
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:head>
    <title>Edit Course | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href="/courses" class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Courses
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Edit Course #{courseId}</h1>

            {#if isLoading}
                <p class="text-gray-500">Loading course data...</p>
            {:else}
                <form onsubmit={handleUpdate} class="flex flex-col gap-5">
                    <div>
                        <label for="edit-teacher-id" class="block text-sm font-medium text-gray-700 mb-1">Teacher ID</label>
                        <input id="edit-teacher-id" type="number" bind:value={teacherId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
                        <input id="edit-name" type="text" bind:value={name} required maxlength="255" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-code" class="block text-sm font-medium text-gray-700 mb-1">Code</label>
                        <input id="edit-code" type="text" bind:value={code} required maxlength="50" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-semester" class="block text-sm font-medium text-gray-700 mb-1">Semester</label>
                        <input id="edit-semester" type="text" bind:value={semester} required maxlength="50" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-max-students" class="block text-sm font-medium text-gray-700 mb-1">Max Students</label>
                        <input id="edit-max-students" type="number" bind:value={maxStudents} required min="1" max="32767" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
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
