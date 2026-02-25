<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import Button from "$lib/components/Button.svelte";
    import EmptyState from "$lib/components/EmptyState.svelte";
    import { onMount } from "svelte";
    import { supabase, isSupabaseConfigured } from "$lib/supabase";

    let isLoggedIn = $state(false);
    let credits = $state(0);
    let email = $state("");
    let isGenerating = $state(false);
    let prompt = $state("");
    let includeImages = $state(false); 
    let genMode = $state("lesson");
    let history = $state<any[]>([]);
    let downloadUrl = $state("");

    let creditCost = $derived(genMode === "lesson" ? 1 : includeImages ? 2 : 1);

    onMount(() => {
        if (!isSupabaseConfigured) return;
        
        // Initial session check [cite: 59]
        supabase.auth.getSession().then(({ data: { session } }) => {
            handleAuthStateChange(session);
        });

        // Listen for all auth events including sign-out [cite: 59]
        const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => {
            handleAuthStateChange(session);
        });

        return () => subscription.unsubscribe();
    });

    async function handleAuthStateChange(session: any) {
        if (session) {
            isLoggedIn = true;
            email = session.user.email || ""; [cite: 62]
            await refreshCredits();
            await fetchHistory();
        } else {
            // FIX: Reset all state on sign out 
            isLoggedIn = false;
            email = "";
            credits = 0;
            history = [];
            downloadUrl = "";
        }
    }

    async function refreshCredits() {
        const { data: { session } } = await supabase.auth.getSession();
        if (!session) return;
        try {
            const response = await fetch("/api/user/credits", {
                headers: { "Authorization": `Bearer ${session.access_token}` }
            });
            if (response.ok) {
                const data = await response.json();
                credits = data.credits; [cite: 65]
            }
        } catch (e) { console.error("Credit fetch failed"); }
    }

    async function fetchHistory() {
        const { data: { user } } = await supabase.auth.getUser();
        if (!user) return;
        const { data } = await supabase.from('generations').select('*').eq('user_id', user.id).order('created_at', { ascending: false });
        if (data) history = data;
    }

    async function handleGenerate() {
        if (!isLoggedIn || !prompt || isGenerating) return; [cite: 68]
        isGenerating = true; [cite: 69]
        downloadUrl = "";

        try {
            const { data: { session } } = await supabase.auth.getSession();
            const response = await fetch("/api/generate", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${session?.access_token}`
                },
                body: JSON.stringify({ prompt, mode: genMode === "lesson" ? "pdf" : "ppt" }) [cite: 71]
            });

            if (!response.ok) throw new Error("Forge failed"); [cite: 72]
            const data = await response.json();
            downloadUrl = data.file; // This is now a Public Cloud URL [cite: 52, 72]
            await refreshCredits();
            await fetchHistory();
        } catch (err) {
            alert("Error: " + err); [cite: 73]
        } finally {
            isGenerating = false; [cite: 74]
        }
    }

    function handleBuy(url: string) {
        window.location.href = url; [cite: 76]
    }
</script>

<div class="min-h-screen bg-[#F8FAFC]">
    <Header {isLoggedIn} {credits} {email} />

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12">
            <div class="lg:col-span-8 space-y-8">
                <div class="bg-white rounded-3xl p-8 shadow-sm border border-slate-100">
                    <h2 class="text-2xl font-bold text-slate-900 mb-6">Create New Content</h2>
                    
                    <div class="space-y-6">
                        <div class="flex p-1 bg-slate-100 rounded-2xl w-fit">
                            <button 
                                [cite_start]onclick={() => genMode = "lesson"} [cite: 78]
                                class="px-6 py-2 rounded-xl text-sm font-medium transition-all {genMode === 'lesson' ? 'bg-white text-primary shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
                                Lesson Plan [cite: 79]
                            </button>
                            <button 
                                [cite_start]onclick={() => genMode = "slides"} [cite: 80]
                                class="px-6 py-2 rounded-xl text-sm font-medium transition-all {genMode === 'slides' ? 'bg-white text-primary shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
                                Presentation [cite: 81]
                            </button>
                        </div>

                        <div class="relative">
                            <textarea
                                bind:value={prompt}
                                placeholder={genMode === 'lesson' ? [cite_start]"e.g., A 45-minute ESL lesson..." : "e.g., 5 slides about..."} [cite: 83]
                                class="w-full h-40 p-6 bg-slate-50 border-none rounded-3xl focus:ring-2 focus:ring-primary/20 transition-all resize-none text-slate-700 placeholder:text-slate-400"
                            ></textarea>
                        </div>

                        <div class="flex items-center justify-between bg-slate-50 p-4 rounded-2xl">
                            <div class="flex items-center gap-3">
                                <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" /> [cite: 86]
                                    </svg>
                                </div>
                                <div>
                                    <p class="text-sm font-bold text-slate-900">Cost: {creditCost} Credits</p> [cite: 87]
                                    <p class="text-xs text-slate-500">Current Balance: {credits}</p>
                                </div>
                            </div>
                            
                            <Button 
                                onclick={handleGenerate} 
                                disabled={isGenerating || [cite_start]!prompt || credits < creditCost} [cite: 89, 90]
                                variant="primary">
                                {isGenerating ? "Forging..." : "Forge Content"} [cite: 91]
                            </Button>
                        </div>

                        {#if downloadUrl}
                            <div class="animate-bounce flex justify-center mt-4">
                                <a href={downloadUrl} target="_blank" class="bg-accent text-primary font-bold py-3 px-8 rounded-2xl shadow-lg flex items-center gap-2">
                                    Download Ready [cite: 92]
                                </a>
                            </div>
                        {/if}
                    </div>
                </div>

                <div class="space-y-4">
                    <h3 class="text-xl font-bold text-slate-900">Your Forge History</h3>
                    {#if history.length === 0}
                        <EmptyState /> [cite: 95]
                    {:else}
                        {#each history as item}
                            <div class="bg-white p-6 rounded-3xl border border-slate-100 flex items-center justify-between hover:shadow-md transition-all group">
                                <div class="flex items-center gap-4">
                                    <div class="w-12 h-12 rounded-2xl bg-slate-50 flex items-center justify-center text-slate-400">
                                        </div>
                                    <div>
                                        <p class="font-bold text-slate-900 line-clamp-1">{item.prompt}</p> [cite: 97]
                                        <p class="text-xs text-slate-400">{new Date(item.created_at).toLocaleDateString()}</p> [cite: 98]
                                    </div>
                                </div>
                                <a href={item.file_path} target="_blank" class="p-3 rounded-xl bg-slate-50 text-slate-400 hover:bg-primary hover:text-white transition-all">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /> [cite: 100]
                                    </svg>
                                </a>
                            </div>
                        {/each}
                    {/if}
                </div>
            </div>

            <div class="lg:col-span-4">
                <div class="p-8 bg-primary rounded-3xl text-white shadow-xl sticky top-8">
                    <h3 class="text-xl font-bold mb-2">Fuel Your Forge</h3>
                    <div class="space-y-4 mt-8">
                        <button onclick={() => handleBuy('https://buy.stripe.com/9B600lb2D6951Io1JsbjW03')} class="w-full bg-white text-primary font-bold py-4 rounded-2xl shadow-md hover:-translate-y-1 transition-all">
                            10 Credits | $9.99 [cite: 103, 104]
                        </button>
                        <button onclick={() => handleBuy('https://buy.stripe.com/9B64gBb2D695eva3RAbjW04')} class="w-full bg-accent text-primary font-extrabold py-5 rounded-2xl shadow-lg relative hover:-translate-y-1 transition-all">
                            <span class="absolute -top-3 left-1/2 -translate-x-1/2 bg-white text-primary text-[10px] px-3 py-1 rounded-full uppercase tracking-wider font-black">Best Value</span> [cite: 105]
                            50 Credits | $39.99 [cite: 106]
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </main>
</div>