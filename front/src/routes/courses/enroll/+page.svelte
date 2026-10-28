<script lang="ts">
    import { onMount } from 'svelte';
    import { PUBLIC_API_URL } from '$env/static/public';

    let studentId = $state<string | number>('');
    let courseId = $state<string | number>('');

    let isSaving = $state(false);
    let errorMessage = $state('');
    let successMessage = $state('');

    onMount(() => {
        if (!localStorage.getItem('paseto_token')) {
            window.location.href = '/login';
        }
    });

    async function handleEnroll(event: SubmitEvent) {
        event.preventDefault();
        isSaving = true;
        errorMessage = '';
        successMessage = '';
        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch(`${PUBLIC_API_URL}/courses/enroll/`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    student_id: Number(studentId),
                    course_id: Number(courseId)
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to enroll student');
            }

            const data = await response.json();
            successMessage = `Student enrolled successfully! Enrollment ID: ${data.data}`;

            studentId = '';
            courseId = '';
        } catch (err: any) {
            errorMessage = err.message || 'Connection error.';
        } finally {
            isSaving = false;
        }
    }
</script>
<svelte:head>
    <title>Enroll Student | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8 flex justify-center">
    <div class="w-full max-w-lg">
        <a href="/courses" class="text-blue-600 hover:underline text-sm font-medium mb-4 inline-block">
            &larr; Back to Courses
        </a>

        <div class="bg-white p-8 rounded-xl shadow-sm border border-gray-200">
            <h1 class="text-2xl font-bold text-gray-800 mb-6">Enroll Student</h1>

            <form onsubmit={handleEnroll} class="flex flex-col gap-5">
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

                {#if successMessage}
                    <div class="bg-green-50 text-green-600 text-sm p-3 rounded-lg text-center">{successMessage}</div>
                {/if}

                <button type="submit" disabled={isSaving} class="w-full bg-green-600 hover:bg-green-700 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg mt-2 transition-colors">
                    {isSaving ? 'Enrolling...' : 'Enroll Student'}
                </button>
            </form>
        </div>
    </div>
</div>
