<script lang="ts">
    import { onMount } from "svelte";
    import { supabase } from "$lib/supabase";

    onMount(async () => {
        // Exchange the code for a session is handled automatically by the Supabase client
        // if it detects the code in the URL. We just need to wait a moment and then redirect.
        const {
            data: { session },
            error,
        } = await supabase.auth.getSession();

        if (error) {
            console.error("Error during auth callback:", error);
            window.location.href =
                "/login?error=" + encodeURIComponent(error.message);
        } else if (session) {
            window.location.href = "/";
        } else {
            // If no session yet, might still be processing or failed silently
            window.location.href = "/login";
        }
    });
</script>

<div class="min-h-screen bg-slate-50 flex items-center justify-center p-4">
    <div class="text-center">
        <div
            class="animate-spin h-10 w-10 text-primary mx-auto mb-4 border-4 border-t-transparent border-primary rounded-full"
        ></div>
        <h2 class="text-xl font-bold text-slate-800">Completing Sign In...</h2>
        <p class="text-slate-500 mt-2">Forging your session, please wait.</p>
    </div>
</div>
