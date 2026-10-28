<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/state';
    import { resolve } from '$app/paths';
    import { PUBLIC_API_URL } from '$env/static/public';

    let enrollmentId = $derived(page.params.id);

    let studentId = $state<string | number>('');
    let courseId = $state<string | number>('');
    let enrolledAt = $state('');
    let status = $state('');

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
            const response = await fetch(`${PUBLIC_API_URL}/enrollments/${enrollmentId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (!response.ok) throw new Error('Enrollment not found');

            const data = await response.json();
            const enrollment = data.data ? data.data : data;

            studentId = enrollment.student_id;
            courseId = enrollment.course_id;
            enrolledAt = new Date(enrollment.enrolled_at).toLocaleDateString();
            status = enrollment.status;
        } catch (err) {
            errorMessage = 'Failed to load enrollment data.';
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
            const response = await fetch(`${PUBLIC_API_URL}/enrollments`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    id: Number(enrollmentId),
                    student_id: Number(studentId),
                    course_id: Number(courseId)
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to update enrollment');
            }

            successMessage = 'Enrollment updated successfully!';
            setTimeout(() => { window.location.href = '/enrollments'; }, 1500);
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:head>
    <title>Edit Enrollment | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href={resolve('/enrollments')} class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Enrollments
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Edit Enrollment #{enrollmentId}</h1>

            {#if isLoading}
                <p class="text-gray-500">Loading enrollment data...</p>
            {:else}
                <form onsubmit={handleUpdate} class="flex flex-col gap-5">
                    <div>
                        <label for="edit-student-id" class="block text-sm font-medium text-gray-700 mb-1">Student ID</label>
                        <input id="edit-student-id" type="number" bind:value={studentId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label for="edit-course-id" class="block text-sm font-medium text-gray-700 mb-1">Course ID</label>
                        <input id="edit-course-id" type="number" bind:value={courseId} required min="1" class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 focus:ring-2 focus:ring-blue-500 focus:outline-none" />
                    </div>

                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Enrolled At</label>
                        <input type="text" value={enrolledAt} disabled class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-200 text-gray-500 cursor-not-allowed" />
                    </div>

                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
                        <input type="text" value={status} disabled class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-200 text-gray-500 cursor-not-allowed uppercase font-semibold" />
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
