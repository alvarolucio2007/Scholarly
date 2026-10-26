<script lang="ts">
    import {resolve} from '$app/paths'
    let email = $state('');
    let password = $state('');
    let errorMessage = $state('');

    async function handleLogin(event: SubmitEvent) {
        event.preventDefault();
        errorMessage = '';

        try {
            const response = await fetch('http://localhost:8080/authentication/token', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            if (!response.ok) {
                errorMessage = 'Invalid email or password.';
                return;
            }

            const data = await response.json();
            localStorage.setItem('paseto_token', data.data.token);
            window.location.href = '/dashboard';
        } catch (error) {
            errorMessage = 'Connection error. Please try again.';
            console.log(error);
        }
    }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100 p-4">
    <div class="w-full max-w-md bg-white p-8 rounded-xl shadow-lg">

        <h1 class="text-2xl font-bold text-center text-gray-800 mb-6">Welcome Back</h1>

        <form onsubmit={handleLogin} class="flex flex-col gap-5">
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                <input
                    type="email"
                    bind:value={email}
                    required
                    placeholder="you@example.com"
                    class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                />
            </div>

            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
                <input
                    type="password"
                    bind:value={password}
                    required
                    placeholder="••••••••"
                    class="w-full p-2.5 rounded-lg border border-gray-300 bg-gray-50 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                />
            </div>

            {#if errorMessage}
                <div class="bg-red-50 text-red-500 text-sm p-3 rounded-lg text-center">
                    {errorMessage}
                </div>
            {/if}

            <button
                type="submit"
                class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2.5 px-4 rounded-lg mt-2 transition-colors duration-200"
            >
                Sign In
            </button>
        </form>

        <div class="mt-6 text-center text-sm text-gray-600">
            Don't have an account?
            <a href={resolve('/register')} class="text-blue-600 hover:text-blue-800 font-semibold hover:underline transition-all">
                Register here
            </a>
        </div>
    </div>
</div>
