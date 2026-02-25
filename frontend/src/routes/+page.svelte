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
    let genMode = $state("lesson");
    let history = $state<any[]>([]);
    let downloadUrl = $state("");

    let creditCost = $derived(genMode === "lesson" ? 1 : 2);

    onMount(() => {
        if (!isSupabaseConfigured) return;
        
        const checkUser = async () => {
            const { data: { session } } = await supabase.auth.getSession();
            if (session) {
                await handleAuthStateChange(session);
            }
        };

        checkUser();

        const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => {
            handleAuthStateChange(session);
        });

        return () => subscription.unsubscribe();
    });

    async function handleAuthStateChange(session: any) {
        if (session) {
            isLoggedIn = true;
            email = session.user.email || "";
            await refreshCredits();
            await fetchHistory();
        } else {
            isLoggedIn = false;
            email = "";
            credits = 0;
            history = [];
        }
    }

    async function refreshCredits() {
        const { data: { session } } = await supabase.auth.getSession();
        if (!session) return;

        try {
            const response = await fetch("/api/user/credits", {
                headers: {
                    "Authorization": `Bearer ${session.access_token}`
                }
            });
            if (response.ok) {
                const data = await response.json();
                credits = data.credits;
            }
        } catch (err) {
            console.error("Failed to fetch credits from Vercel function");
        }
    }

    async function fetchHistory() {
        const { data: { user } } = await supabase.auth.getUser();
        if (!user) return;

        const { data, error } = await supabase
            .from('generations')
            .select('*')
            .eq('user_id', user.id)
            .order('created_at', { ascending: false });
        
        if (!error && data) {
            history = data;
        }
    }

    async function handleGenerate() {
        if (!isLoggedIn || !prompt || isGenerating) return;
        
        isGenerating = true;
        downloadUrl = "";

        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) throw new Error("Please log in again.");

            const response = await fetch("/api/generate", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${session.access_token}`
                },
                body: JSON.stringify({
                    prompt: prompt,
                    mode: genMode === "lesson" ? "pdf" : "ppt"
                })
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || "Server error");
            }

            const data = await response.json();
            
            if (data.file) {
                downloadUrl = data.file;
                await refreshCredits();
                await fetchHistory();
            }
        } catch (err: any) {
            console.error("Generation error:", err);
            alert("Forge Failed: " + err.message);
        } finally {
            isGenerating = false;
        }
    }

    function handleBuy(url: string) {
        window.location.href = url;
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
                                onclick={() => genMode = "lesson"}
                                class="px-6 py-2 rounded-xl text-sm font-medium transition-all {genMode === 'lesson' ? 'bg-white text-primary shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
                                Lesson Plan
                            </button>
                            <button 
                                onclick={() => genMode = "slides"}
                                class="px-6 py-2 rounded-xl text-sm font-medium transition-all {genMode === 'slides' ? 'bg-white text-primary shadow-sm' : 'text-slate-500 hover:text-slate-700'}">
                                Presentation
                            </button>
                        </div>

                        <div class="relative">
                            <textarea
                                bind:value={prompt}
                                placeholder={genMode === 'lesson' ? "Describe your lesson..." : "Describe your slides..."}
                                class="w-full h-40 p-6 bg-slate-50 border-none rounded-3xl focus:ring-2 focus:ring-primary/20 transition-all resize-none text-slate-700"
                            ></textarea>
                        </div>

                        <div class="flex items-center justify-between bg-slate-50 p-4 rounded-2xl">
                            <div class="flex items-center gap-3">
                                <div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary font-bold">
                                    !
                                </div>
                                <div>
                                    <p class="text-sm font-bold text-slate-900">Cost: {creditCost} Credits</p>
                                    <p class="text-xs text-slate-500">Balance: {credits}</p>
                                </div>
                            </div>
                            
                            <Button 
                                onclick={handleGenerate} 
                                disabled={isGenerating || !prompt || credits < creditCost}
                                variant="primary">
                                {isGenerating ? "Forging..." : "Forge Content"}
                            </Button>
                        </div>

                        {#if downloadUrl}
                            <div class="flex justify-center mt-4">
                                <a href={downloadUrl} target="_blank" class="bg-emerald-500 text-white font-bold py-3 px-8 rounded-2xl shadow-lg flex items-center gap-2 hover:bg-emerald-600 transition-colors">
                                    Download Ready
                                </a>
                            </div>
                        {/if}
                    </div>
                </div>

                <div class="space-y-4">
                    <h3 class="text-xl font-bold text-slate-900">History</h3>
                    {#if history.length === 0}
                        <EmptyState />
                    {:else}
                        {#each history as item}
                            <div class="bg-white p-6 rounded-3xl border border-slate-100 flex items-center justify-between shadow-sm">
                                <div>
                                    <p class="font-bold text-slate-900">{item.prompt.substring(0, 50)}...</p>
                                    <p class="text-xs text-slate-400">{new Date(item.created_at).toLocaleDateString()}</p>
                                </div>
                                <a href={item.file_path} target="_blank" class="text-primary font-bold hover:underline">
                                    Download
                                </a>
                            </div>
                        {/each}
                    {/if}
                </div>
            </div>

            <div class="lg:col-span-4">
                <div class="p-8 bg-primary rounded-3xl text-white shadow-xl sticky top-8">
                    <h3 class="text-xl font-bold mb-2">Buy Credits</h3>
                    <div class="space-y-4 mt-6">
                        <button onclick={() => handleBuy('https://buy.stripe.com/9B600lb2D6951Io1JsbjW03')} class="w-full bg-white text-primary font-bold py-4 rounded-2xl hover:bg-slate-100">
                            10 Credits | $9.99
                        </button>
                        <button onclick={() => handleBuy('https://buy.stripe.com/9B64gBb2D695eva3RAbjW04')} class="w-full bg-yellow-400 text-primary font-bold py-4 rounded-2xl hover:bg-yellow-500">
                            50 Credits | $39.99
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </main>
</div>