<script lang="ts">
    import { onMount } from 'svelte';
    import { resolve } from '$app/paths';

    // Variáveis de estado
    let users = $state<any[]>([]);
    let isLoading = $state(true);
    let errorMessage = $state('');

    // Busca os usuários no backend
    async function fetchUsers() {
        isLoading = true;
        errorMessage = '';
        const token = localStorage.getItem('paseto_token');

        try {
            const response = await fetch('http://localhost:8080/users', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.status === 401) {
                // Token expirado ou inválido
                localStorage.removeItem('paseto_token');
                window.location.href = '/login';
                return;
            }

            if (!response.ok) {
                throw new Error('Failed to fetch users');
            }

            const data = await response.json();
            // Lida com o envelope {"data": [...]} ou um array puro
            users = data.data ? data.data : data;

        } catch (error) {
            errorMessage = 'Could not load users. Please try again.';
        } finally {
            isLoading = false;
        }
    }

    // Deleta um usuário
    async function deleteUser(id: number) {
        if (!confirm('Are you sure you want to delete this user?')) return;

        const token = localStorage.getItem('paseto_token');
        try {
            const response = await fetch(`http://localhost:8080/users/${id}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (response.ok) {
                // Remove o usuário da lista na tela sem precisar recarregar a página
                users = users.filter(u => u.id !== id);
            } else {
                alert('Failed to delete user.');
            }
        } catch (error) {
            alert('Connection error.');
        }
    }

    // Executa ao carregar a tela
    onMount(() => {
        if (!localStorage.getItem('paseto_token')) {
            window.location.href = '/login';
            return;
        }
        fetchUsers();
    });
</script>

<svelte:head>
    <title>Users | Scholarly</title>
</svelte:head>

<div class="min-h-screen bg-gray-100 p-8">
    <header class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div>
            <a href={resolve('/dashboard')} class="text-blue-600 hover:underline text-sm font-medium mb-2 inline-block">
                &larr; Back to Dashboard
            </a>
            <h1 class="text-3xl font-bold text-gray-800">Users Management</h1>
            <p class="text-gray-500">View, edit, and manage system users</p>
        </div>
    </header>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">

        {#if isLoading}
            <div class="p-10 text-center text-gray-500 font-medium">Loading users...</div>
        {:else if errorMessage}
            <div class="p-10 text-center text-red-500 font-medium">{errorMessage}</div>
        {:else if users.length === 0}
            <div class="p-10 text-center text-gray-500 font-medium">No users found in the system.</div>
        {:else}
            <!-- Tabela de Usuários -->
            <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-gray-50 border-b border-gray-200 text-gray-600 text-sm">
                            <th class="p-4 font-semibold">ID</th>
                            <th class="p-4 font-semibold">Name</th>
                            <th class="p-4 font-semibold">Email</th>
                            <th class="p-4 font-semibold">CPF</th>
                            <th class="p-4 font-semibold text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100">
                        {#each users as user}
                            <tr class="hover:bg-gray-50 transition-colors">
                                <td class="p-4 text-gray-500 font-mono text-sm">#{user.id}</td>
                                <td class="p-4 text-gray-800 font-medium">{user.name}</td>
                                <td class="p-4 text-gray-600">{user.email}</td>
                                <td class="p-4 text-gray-600">{user.cpf || 'Not provided'}</td>
                                <td class="p-4 text-right flex justify-end gap-2">
                                    <a
                                        href={resolve('/users/[id]', { id: user.id.toString() })}
                                        class="text-blue-600 hover:bg-blue-50 px-3 py-1 rounded-md text-sm font-medium transition-colors"
                                    >
                                        Edit
                                    </a>
                                    <button
                                        onclick={() => deleteUser(user.id)}
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
