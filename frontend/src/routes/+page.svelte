<script lang="ts">
    import Header from "$lib/components/Header.svelte";
    import Button from "$lib/components/Button.svelte";
    import EmptyState from "$lib/components/EmptyState.svelte";
    import { onMount } from "svelte";
    import { supabase, isSupabaseConfigured } from "$lib/supabase";

    // --- State Management (Svelte 5 Runes) ---
    let isLoggedIn = $state(false);
    let credits = $state(0);
    let email = $state("");
    let isGenerating = $state(false);
    let prompt = $state("");
    let includeImages = $state(false);
    let genMode = $state("ppt"); // UI state: 'lesson' or 'ppt'
    let history = $state<any[]>([]);

    // --- Derived Logic ---
    let creditCost = $derived(genMode === "lesson" ? 1 : includeImages ? 2 : 1);

    onMount(() => {
        if (!isSupabaseConfigured) {
            console.warn("Supabase is not configured. Check environment variables.");
            return;
        }

        const setupAuth = async () => {
            const { data: { user } } = await supabase.auth.getUser();
            if (user) {
                isLoggedIn = true;
                email = user.email || "";
                await refreshCredits(user.id);
                await fetchHistory(user.id);
            }

            const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => {
                handleAuthStateChange(session);
            });

            return subscription;
        };

        const subscriptionPromise = setupAuth();
        return () => {
            subscriptionPromise.then((sub) => sub?.unsubscribe());
        };
    });

    // --- Helper Functions ---
    async function refreshCredits(userId: string) {
        const { data, error } = await supabase
            .from("users")
            .select("credit_balance")
            .eq("id", userId)
            .single();

        if (data) {
            credits = data.credit_balance;
        } else if (error) {
            console.warn("UI sync error (credits):", error.message);
        }
    }

    async function fetchHistory(userId: string) {
    console.log("🛠️ fetchHistory started for user:", userId);
    
    const { data, error } = await supabase
        .from("generations")
        .select("*")
        .eq("user_id", userId)
        .order("created_at", { ascending: false });

    if (error) {
        console.error("❌ Supabase Error fetching history:", error.message);
        return;
    }

    console.log("📦 Data received from DB:", data);

    if (data) {
        history = data.map((item: any) => ({
            id: item.id,
            title: item.prompt.substring(0, 30) + (item.prompt.length > 30 ? "..." : ""),
            type: item.file_path?.includes("lesson") ? "lesson" : "ppt",
            date: new Date(item.created_at).toLocaleDateString(),
            status: item.status,
            hasPPT: item.file_path?.includes("ppt") || item.status === "completed",
            hasLesson: item.file_path?.includes("lesson") || item.status === "completed",
        }));
        console.log("✨ History state updated. Items count:", history.length);
    }
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
            email = "";
            history = [];
        }
    }

    async function handleSignOut() {
        await supabase.auth.signOut();
    }

    // --- Core Action: Generate ---
    async function handleGenerate() {
        if (!isLoggedIn || !prompt) return;
        
        if (credits < creditCost) {
            alert("Insufficient credits. Please top up your forge!");
            return;
        }

        isGenerating = true;
        try {
            const { data: { session } } = await supabase.auth.getSession();
            if (!session) throw new Error("Authentication session expired.");

            const response = await fetch("http://localhost:8080/api/generate", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${session.access_token}`,
                },
                body: JSON.stringify({ 
                    prompt, 
                    mode: genMode === 'ppt' ? 'pptx' : 'lesson',
                    includeImages 
                }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || "The forge failed to strike.");
            }

            const { data: { user } } = await supabase.auth.getUser();
            if (user) {
                await refreshCredits(user.id);
                await fetchHistory(user.id);
            }
            
            prompt = ""; 
            alert("Magic Forged! Check your Recent Activity.");
            
        } catch (error: any) {
            console.error("Forge Error:", error);
            alert(`Error: ${error.message}`);
        } finally {
            isGenerating = false;
        }
    }
</script>

<svelte:head>
    <title>Lesson Forge | AI Lesson Planner</title>
</svelte:head>

<Header title="Lesson Forge" {email} {credits} {isLoggedIn} onSignOut={handleSignOut} />

<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div class="text-center mb-16">
        <h2 class="text-4xl md:text-5xl font-extrabold text-slate-900 tracking-tight mb-4">
            Helping Forge Future <span class="text-primary italic">Minds</span>
        </h2>
        <p class="text-lg text-slate-500 max-w-2xl mx-auto">
            Generate high-quality presentations and lesson plans in seconds.
        </p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-10">
        <div class="lg:col-span-8 space-y-8">
            <div class="card-premium p-8 relative overflow-hidden bg-white border rounded-3xl shadow-sm">
                <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
                    <h3 class="text-xl font-bold text-slate-800 flex items-center gap-2">Forge Magic</h3>

                    <div class="flex p-1 bg-slate-100 rounded-xl relative z-10">
                        <button
                            class="px-4 py-2 text-sm font-bold rounded-lg transition-all {genMode === 'lesson' ? 'bg-white text-primary shadow-sm' : 'text-slate-500'}"
                            onclick={() => (genMode = "lesson")}
                        >Lesson Plan</button>
                        <button
                            class="px-4 py-2 text-sm font-bold rounded-lg transition-all {genMode === 'ppt' ? 'bg-white text-primary shadow-sm' : 'text-slate-500'}"
                            onclick={() => (genMode = "ppt")}
                        >Presentation</button>
                    </div>
                </div>

                <div class="space-y-6">
                    <div>
                        <label for="prompt" class="block text-sm font-bold text-slate-700 mb-2 uppercase tracking-wide">
                            {genMode === "lesson" ? "Topic & Objectives" : "Presentation Context"}
                        </label>
                        <textarea
                            id="prompt"
                            rows="5"
                            bind:value={prompt}
                            class="w-full px-5 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:ring-4 focus:ring-primary/10 focus:border-primary focus:bg-white transition-all text-slate-800"
                            placeholder={genMode === "lesson" ? "E.g., A lesson for beginner adults about grocery shopping..." : "E.g., Presentation about remote team management..."}
                        ></textarea>
                    </div>

                    {#if genMode === "ppt"}
                        <div class="p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 flex items-center justify-between group hover:border-accent transition-colors">
                            <div class="flex items-center gap-4">
                                <div class="w-10 h-10 rounded-xl bg-white flex items-center justify-center text-accent shadow-sm">
                                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                                </div>
                                <div>
                                    <h4 class="text-sm font-bold text-slate-800">Include AI Images</h4>
                                    <p class="text-xs text-slate-500">Add 4-6 high-quality images</p>
                                </div>
                            </div>
                            <label class="relative inline-flex items-center cursor-pointer">
                                <input type="checkbox" bind:checked={includeImages} class="sr-only peer" />
                                <div class="w-11 h-6 bg-slate-200 rounded-full peer peer-checked:after:translate-x-full peer-checked:bg-primary after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all"></div>
                            </label>
                        </div>
                    {/if}

                    <div class="flex items-center justify-between pt-4 border-t border-slate-100">
                        <div class="text-sm text-slate-500">
                            {#if !isLoggedIn}
                                <span class="text-amber-600 font-bold text-[10px] uppercase tracking-wider bg-amber-50 px-3 py-1 rounded-full">Sign in to start forging</span>
                            {/if}
                        </div>
                        <Button
                            text={isGenerating ? "Forging..." : `Generate (${creditCost} Credit${creditCost > 1 ? "s" : ""})`}
                            isLoading={isGenerating}
                            disabled={!isLoggedIn || !prompt}
                            onclick={handleGenerate}
                        />
                    </div>
                </div>
            </div>

            <div class="card-premium overflow-hidden bg-white border rounded-3xl shadow-sm">
                <div class="px-8 py-6 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
                    <h3 class="text-lg font-bold text-slate-800 flex items-center gap-2">Recent Activity</h3>
                </div>
                <div class="divide-y divide-slate-50">
                    {#if history.length === 0}
                        <div class="px-8 py-12"><EmptyState message="Your forged materials will appear here." /></div>
                    {:else}
                        {#each history as item}
                            <div class="px-8 py-6 flex items-center justify-between hover:bg-slate-50/50 transition-colors group">
                                <div class="flex items-center gap-4">
                                    <div class="w-10 h-10 rounded-xl flex items-center justify-center {item.type === 'lesson' ? 'bg-emerald-50 text-emerald-600' : 'bg-blue-50 text-blue-600'}">
                                        {#if item.type === "lesson"}
                                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" /></svg>
                                        {:else}
                                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" /></svg>
                                        {/if}
                                    </div>
                                    <div>
                                        <h4 class="text-sm font-bold text-slate-800 line-clamp-1">{item.title}</h4>
                                        <p class="text-xs text-slate-500 uppercase font-semibold">{item.type} • {item.date}</p>
                                    </div>
                                </div>
                                <div class="flex items-center gap-3">
                                    <button class="p-2 text-slate-400 hover:text-primary transition-colors">
                                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                                    </button>
                                </div>
                            </div>
                        {/each}
                    {/if}
                </div>
            </div>
        </div>

        <div class="lg:col-span-4 space-y-8">
            <div class="p-8 bg-primary rounded-3xl text-white shadow-xl">
                <h3 class="text-xl font-bold mb-2">Fuel Your Forge</h3>
                <p class="text-white/70 text-sm mb-8">Unlock unlimited creativity with credit bundles.</p>
                <div class="space-y-4">
                    <button class="w-full bg-white text-primary font-bold py-4 rounded-2xl shadow-xl hover:-translate-y-1 transition-all">
                        10 Credits | $9.99
                    </button>
                    <button class="w-full bg-accent text-primary font-extrabold py-5 rounded-2xl shadow-2xl relative">
                        <span class="absolute -top-3 left-1/2 -translate-x-1/2 bg-secondary text-white text-[10px] px-3 py-1 rounded-full">Popular</span>
                        50 Credits | $39.99
                    </button>
                </div>
            </div>
        </div>
    </div>
</main>