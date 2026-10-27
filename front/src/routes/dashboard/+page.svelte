<script lang="ts">
    import { onMount } from 'svelte';
    import { resolve } from '$app/paths';

    let userCount = $state(0);
    let courseCount = $state(0);
    let studentCount = $state(0);
    let isLoading = $state(true);

    async function fetchStats() {
        const token = localStorage.getItem('paseto_token');
        const headers = { 'Authorization': `Bearer ${token}` };

        try {
            const [usersRes, coursesRes, studentsRes] = await Promise.all([
                fetch('http://localhost:8080/users', { headers }),
                fetch('http://localhost:8080/courses', { headers }),
                fetch('http://localhost:8080/students', { headers })
            ]);

            if (usersRes.ok) userCount = (await usersRes.json()).data.length;
            if (coursesRes.ok) courseCount = (await coursesRes.json()).data.length;
            if (studentsRes.ok) studentCount = (await studentsRes.json()).data.length;

        } catch (error) {
            console.error("Error while fetching statistics", error);
        } finally {
            isLoading = false;
        }
    }

    onMount(() => {
        if (!localStorage.getItem('paseto_token')) {
            window.location.href = '/login';
            return;
        }
        fetchStats();
    });
</script>

<div class="min-h-screen bg-gray-100 p-8">
    <header class="flex justify-between items-center mb-10">
        <div>
            <h1 class="text-3xl font-bold text-gray-800">Dashboard</h1>
            <p class="text-gray-500">Welcome to Scholarly!</p>
        </div>
        <button
            onclick={() => { localStorage.removeItem('paseto_token'); window.location.href = '/login'; }}
            class="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg font-semibold transition-colors"
        >
            Sair
        </button>
    </header>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
        <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
            <h3 class="text-gray-500 text-sm font-semibold">User Count</h3>
            <p class="text-4xl font-bold text-blue-600 mt-2">{isLoading ? '...' : userCount}</p>
        </div>
        <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
            <h3 class="text-gray-500 text-sm font-semibold">Enrolled Students</h3>
            <p class="text-4xl font-bold text-green-600 mt-2">{isLoading ? '...' : studentCount}</p>
        </div>
        <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
            <h3 class="text-gray-500 text-sm font-semibold">Active courses</h3>
            <p class="text-4xl font-bold text-purple-600 mt-2">{isLoading ? '...' : courseCount}</p>
        </div>
    </div>

    <h2 class="text-xl font-bold text-gray-800 mb-4">Quick access</h2>
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <a href={resolve('/users')} class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 hover:shadow-md hover:border-blue-400 transition-all flex flex-col items-center justify-center text-center">
            <span class="text-3xl mb-2">👥</span>
            <span class="font-semibold text-gray-700">Users</span>
        </a>
        <a href={resolve('/courses')} class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 hover:shadow-md hover:border-blue-400 transition-all flex flex-col items-center justify-center text-center">
            <span class="text-3xl mb-2">📚</span>
            <span class="font-semibold text-gray-700">Courses</span>
        </a>
        <a href={resolve('/students')} class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 hover:shadow-md hover:border-blue-400 transition-all flex flex-col items-center justify-center text-center">
            <span class="text-3xl mb-2">🎓</span>
            <span class="font-semibold text-gray-700">Students</span>
        </a>
        <a href={resolve('/grades')} class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 hover:shadow-md hover:border-blue-400 transition-all flex flex-col items-center justify-center text-center">
            <span class="text-3xl mb-2">📝</span>
            <span class="font-semibold text-gray-700">Grades</span>
        </a>
    </div>
</div>
