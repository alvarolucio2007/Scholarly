<script lang="ts">
    let email=$state('');
    let password = $state('');
    let errorMessage=$state('');
    async function handleLogin(event:SubmitEvent){
      event.preventDefault();
      errorMessage='';
        const response=await fetch('http://localhost:8080/authentication/token',{
          method:'POST',
          headers:{
            'Content-Type':'application/json'
          },
          body:JSON.stringify({email,password})
        });
        if (!response.ok){
          errorMessage='Invalid email or passowrd.';
          return;
        }
        const data=await response.json();
        localStorage.setItem('paseto_token',data.data.token);
        window.location.href='/dashboard'
    }
</script>
<div class="flex flex-col items-center justify-center h-screen bg-gray900 text-white">
    <div class="bg-gray-800 p-8 rounded-lg shadow-md w-96">
        <form onsubmit={handleLogin} class="flex flex-col gap-4">
                    <div>
                        <label class="block text-sm mb-1">Email</label>
                        <input
                            type="email"
                            bind:value={email}
                            required
                            class="w-full p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
                        />
                    </div>
                    <div>
                        <label class="block text-sm mb-1">Password</label>
                        <input
                            type="password"
                            bind:value={password}
                            required
                            class="w-full p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
                        />
                    </div>
                    {#if errorMessage}
                        <p class="text-red-500 text-sm text-center">{errorMessage}</p>
                    {/if}

                    <button
                        type="submit"
                        class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded mt-2 transition-colors"
                    >
                        Entrar
                    </button>
                </form>
    </div>
</div>
