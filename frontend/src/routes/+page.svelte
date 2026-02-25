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
        const setupAuth = async () => {
            const { data: { user } } = await supabase.auth.getUser();
            if (user) {
                isLoggedIn = true;
                email = user.email || "";
                await refreshCredits(user.id);
                await fetchHistory(user.id);
            }
            return supabase.auth.onAuthStateChange((_event, session) => {
                handleAuthStateChange(session);
            });
        };
        setupAuth();
    });

    async function handleGenerate() {
        if (!isLoggedIn || !prompt) return;
        isGenerating = true;
        downloadUrl = "";

        try {
            const { data: { session } } = await supabase.auth.getSession();
            const response = await fetch("http://localhost:8080/api/generate", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${session?.access_token}`,
                },
                body: JSON.stringify({ prompt, mode: genMode, includeImages }),
            });

            if (!response.ok) throw new Error("Forge failed.");
            const result = await response.json();

            if (result.file) {
                const path = result.file;
                if (path.startsWith('http')) {
                    downloadUrl = path;
                } else {
                    const cleanPath = path.startsWith('/') ? path : `/${path}`;
                    downloadUrl = `http://localhost:8080${cleanPath}`;
                }

                const newActivity = {
                    id: crypto.randomUUID(),
                    prompt: prompt,
                    file_path: downloadUrl,
                    created_at: new Date().toISOString()
                };
                history = [newActivity, ...history];
            }

            const { data: { user } } = await supabase.auth.getUser();
            if (user) {
                await refreshCredits(user.id);
            }
        } catch (error) {
            console.error("Forge failed", error);
        } finally {
            isGenerating = false;
        }
    }

    async function refreshCredits(userId: string) {
        const { data } = await supabase.from("users").select("credit_balance").eq("id", userId).single();
        if (data) credits = data.credit_balance;
    }

    async function fetchHistory(userId: string) {
        const { data } = await supabase.from("generations").select("*").eq("user_id", userId).order("created_at", { ascending: false });
        if (data) history = data;
    }

    async function handleAuthStateChange(session: any) {
        if (session) {
            isLoggedIn = true;
            email = session.user.email;
            await refreshCredits(session.user.id);
            await fetchHistory(session.user.id);
        } else {
            isLoggedIn = false;
            credits = 0;
            history = [];
        }
    }

    async function handleSignOut() { await supabase.auth.signOut(); }

    function handleBuy(url: string) {
        window.location.href = url;
    }
</script>

<Header title="Lesson Forge" {email} {credits} {isLoggedIn} onSignOut={handleSignOut} />

<main class="max-w-7xl mx-auto px-4 py-12">
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-10">
        <div class="lg:col-span-8 space-y-8">
            <div class="bg-white p-8 rounded-3xl border border-slate-200 shadow-sm">
                <div class="flex items-center justify-between mb-8">
                    <h3 class="text-xl font-bold text-slate-900">Forge Magic</h3>
                    <div class="flex p-1 bg-slate-100 rounded-xl">
                        <button class="px-4 py-2 text-sm font-bold rounded-lg {genMode === 'lesson' ? 'bg-white text-primary shadow-sm' : 'text-slate-500'}" onclick={() => genMode = 'lesson'}>Lesson Plan</button>
                        <button class="px-4 py-2 text-sm font-bold rounded-lg {genMode === 'ppt' ? 'bg-white text-primary shadow-sm' : 'text-slate-500'}" onclick={() => genMode = 'ppt'}>Presentation</button>
                    </div>
                </div>

                <textarea bind:value={prompt} class="w-full h-40 p-6 bg-slate-50 border border-slate-200 rounded-2xl outline-none text-slate-700" placeholder="A B1 level lesson about..."></textarea>
                
                <div class="mt-6 flex items-center justify-between">
                    <div class="flex items-center gap-6">
                        {#if genMode === 'ppt'}
                            <label class="flex items-center gap-2 cursor-pointer">
                                <input type="checkbox" bind:checked={includeImages} class="w-4 h-4 rounded text-primary focus:ring-primary" />
                                <span class="text-sm font-medium text-slate-600">Include Images (+1 Credit)</span>
                            </label>
                        {/if}
                        <div class="text-sm font-medium text-slate-500">Cost: <span class="text-primary font-bold">{creditCost} Credits</span></div>
                    </div>

                    <div class="flex items-center gap-4">
                        {#if downloadUrl}
                            <a href={downloadUrl} target="_blank" rel="noopener noreferrer" download class="px-6 py-3 bg-primary text-white font-bold rounded-xl hover:opacity-90 transition-all shadow-md flex items-center gap-2">
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                                Download Now
                            </a>
                        {/if}
                        <Button text={isGenerating ? "Forging..." : "Forge Magic"} isLoading={isGenerating} disabled={!isLoggedIn || !prompt} onclick={handleGenerate} />
                    </div>
                </div>
            </div>

            <div class="bg-white rounded-3xl border border-slate-200 shadow-sm overflow-hidden">
                <div class="px-8 py-6 border-b border-slate-100"><h3 class="font-bold text-slate-900">Recent Activity</h3></div>
                <div class="divide-y divide-slate-50">
                    {#if history.length === 0}
                        <div class="p-12 text-center"><EmptyState message="Your forge is cold." /></div>
                    {:else}
                        {#each history as item}
                            <div class="px-8 py-6 flex items-center justify-between hover:bg-slate-50 transition-colors">
                                <div class="flex items-center gap-4">
                                    <div class="w-10 h-10 rounded-xl flex items-center justify-center {item.file_path.endsWith('.pptx') ? 'bg-orange-100 text-orange-600' : 'bg-slate-100 text-slate-500'}">
                                        {#if item.file_path.endsWith('.pptx')}
                                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6" />
                                            </svg>
                                        {:else}
                                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                            </svg>
                                        {/if}
                                    </div>
                                    <div class="font-bold text-slate-900">{item.prompt.substring(0, 40)}...</div>
                                </div>
                                <a href={item.file_path.startsWith('http') ? item.file_path : `http://localhost:8080${item.file_path.startsWith('/') ? item.file_path : '/' + item.file_path}`} target="_blank" rel="noopener noreferrer" download class="p-2 text-slate-400 hover:text-primary transition-colors">
                                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                                </a>
                            </div>
                        {/each}
                    {/if}
                </div>
            </div>
        </div>

        <div class="lg:col-span-4">
            <div class="p-8 bg-primary rounded-3xl text-white shadow-xl">
                <h3 class="text-xl font-bold mb-2">Fuel Your Forge</h3>
                <p class="text-white/70 text-sm mb-8">Unlock unlimited creativity.</p>
                <div class="space-y-4">
                    <button onclick={() => handleBuy('https://buy.stripe.com/9B600lb2D6951Io1JsbjW03')} class="w-full bg-white text-primary font-bold py-4 rounded-2xl shadow-md hover:-translate-y-1 transition-all">
                        10 Credits | $9.99
                    </button>
                    <button onclick={() => handleBuy('https://buy.stripe.com/9B64gBb2D695eva3RAbjW04')} class="w-full bg-accent text-primary font-extrabold py-5 rounded-2xl shadow-lg relative hover:-translate-y-1 transition-all">
                        <span class="absolute -top-3 left-1/2 -translate-x-1/2 bg-secondary text-white text-[10px] px-3 py-1 rounded-full uppercase tracking-wider font-bold">Popular</span>
                        50 Credits | $39.99
                    </button>
                </div>
            </div>
        </div>
    </div>
</main>