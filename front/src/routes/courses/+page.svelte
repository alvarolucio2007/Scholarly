<script lang="ts">
    import { onMount } from 'svelte';
    import { PUBLIC_API_URL } from '$env/static/public';

    let courses = $state<any[]>([]);
    let isLoading = $state(true);
    let errorMessage = $state('');

    async function fetchCourses() {
        isLoading = true;
        errorMessage = '';
        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch(`${PUBLIC_API_URL}/courses`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (response.status === 401) {
                localStorage.removeItem('paseto_token');
                window.location.href = '/login';
                return;
            }

            if (!response.ok) throw new Error('Failed to fetch courses');

            const data = await response.json();
            courses = data.data ? data.data : data;

        } catch (error) {
            errorMessage = 'Could not load courses.';
        } finally {
            isLoading = false;
        }
    }

    async function deleteCourse(id: number) {
        if (!confirm('Are you sure you want to delete this course?')) return;

        const token = localStorage.getItem('paseto_token');
        try {
            const response = await fetch(`${PUBLIC_API_URL}/courses/${id}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (response.ok) {
                courses = courses.filter(c => c.id !== id);
            } else {
                alert('Failed to delete course.');
            }
        } catch (error) {
            alert('Connection error.');
        }
    }

    onMount(() => {
        if (!localStorage.getItem('paseto_token')) {
            window.location.href = '/login';
            return;
        }
        fetchCourses();
    });
</script>
<svelte:head>
    <title>Courses | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8">
    <header class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
            <a href="/dashboard" class="text-blue-600 hover:underline text-sm font-medium mb-2 inline-block">
                &larr; Back to Dashboard
            </a>
            <h1 class="text-3xl font-bold text-gray-800">Courses Management</h1>
            <p class="text-gray-500">View and manage available courses</p>
        </div>

        <div class="flex gap-2">
            <a href="/courses/enroll" class="bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-6 rounded-lg transition-colors inline-block">
                Enroll Student
            </a>
            <a href="/courses/create" class="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-6 rounded-lg transition-colors inline-block">
                + New Course
            </a>
        </div>
    </header>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {#if isLoading}
            <div class="p-10 text-center text-gray-500 font-medium">Loading courses...</div>
        {:else if errorMessage}
            <div class="p-10 text-center text-red-500 font-medium">{errorMessage}</div>
        {:else if courses.length === 0}
            <div class="p-10 text-center text-gray-500 font-medium">No courses found.</div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-gray-50 border-b border-gray-200 text-gray-600 text-sm">
                            <th class="p-4 font-semibold">ID</th>
                            <th class="p-4 font-semibold">Name</th>
                            <th class="p-4 font-semibold">Code</th>
                            <th class="p-4 font-semibold">Semester</th>
                            <th class="p-4 font-semibold">Teacher ID</th>
                            <th class="p-4 font-semibold">Max Students</th>
                            <th class="p-4 font-semibold text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100">
                        {#each courses as course}
                            <tr class="hover:bg-gray-50 transition-colors">
                                <td class="p-4 text-gray-500 font-mono text-sm">#{course.id}</td>
                                <td class="p-4 text-gray-800 font-medium">{course.name}</td>
                                <td class="p-4 text-gray-600">{course.code}</td>
                                <td class="p-4 text-gray-600">{course.semester}</td>
                                <td class="p-4 text-gray-500 font-mono text-sm">#{course.teacher_id}</td>
                                <td class="p-4 text-gray-600">{course.max_students}</td>
                                <td class="p-4 text-right flex justify-end gap-2">
                                    <a
                                        href={`/courses/${course.id}`}
                                        class="text-blue-600 hover:bg-blue-50 px-3 py-1 rounded-md text-sm font-medium transition-colors"
                                    >
                                        Edit
                                    </a>
                                    <button
                                        onclick={() => deleteCourse(course.id)}
                                        class="text-red-600 hover:bg-red-50 px-3 py-1 rounded-md text-sm font-medium transition-colors"
                                    >
                                        Delete
                                    </button>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
</div>
