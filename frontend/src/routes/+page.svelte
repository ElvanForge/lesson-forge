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
    let genMode = $state("ppt"); // 'lesson' or 'ppt'

    let creditCost = $derived(genMode === "lesson" ? 1 : includeImages ? 2 : 1);

    onMount(() => {
        if (!isSupabaseConfigured) {
            console.warn(
                "Supabase is not configured. Please set environment variables.",
            );
            return;
        }

        const setupAuth = async () => {
            const {
                data: { session },
            } = await supabase.auth.getSession();
            if (session) {
                handleAuthStateChange(session);
            }

            const {
                data: { subscription },
            } = supabase.auth.onAuthStateChange((_event, session) => {
                handleAuthStateChange(session);
            });

            return subscription;
        };

        const subscriptionPromise = setupAuth();

        return () => {
            subscriptionPromise.then((sub) => sub.unsubscribe());
        };
    });

    async function handleAuthStateChange(session: any) {
        if (session) {
            isLoggedIn = true;
            email = session.user.email;
            // Fetch credits from Supabase
            const { data, error } = await supabase
                .from("users")
                .select("credit_balance")
                .eq("id", session.user.id)
                .single();

            if (data) {
                credits = data.credit_balance;
            } else {
                console.warn(
                    "User profile not found yet. It may be initializing...",
                );
            }
            // Fetch history when logged in
            fetchHistory(session.user.id);
        } else {
            isLoggedIn = false;
            credits = 0;
            email = "";
            history = [];
        }
    }

    async function fetchHistory(userId: string) {
        const { data, error } = await supabase
            .from("generations")
            .select("*")
            .eq("user_id", userId)
            .order("created_at", { ascending: false });

        if (data) {
            history = data.map((item: any) => ({
                id: item.id,
                title: item.prompt.substring(0, 30) + "...",
                type: item.file_path.includes("lesson") ? "lesson" : "ppt", // Simple heuristic
                date: new Date(item.created_at).toLocaleDateString(),
                status: item.status,
                hasPPT:
                    item.file_path.includes("ppt") ||
                    item.status === "completed",
                hasLesson:
                    item.file_path.includes("lesson") ||
                    item.status === "completed",
            }));
        }
    }

    async function handleSignOut() {
        await supabase.auth.signOut();
    }

    let history = $state<any[]>([]);

    async function handleGenerate() {
        if (!isLoggedIn) return;
        if (!prompt) return;
        if (credits < creditCost) {
            alert("Insufficient credits. Please top up.");
            return;
        }

        isGenerating = true;
        try {
            const {
                data: { session },
            } = await supabase.auth.getSession();
            if (!session) throw new Error("No active session");

            const response = await fetch("http://localhost:8080/api/generate", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${session.access_token}`,
                },
                body: JSON.stringify({ prompt, mode: genMode, includeImages }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || "Generation failed");
            }

            const result = await response.json();

            // Re-fetch credits and history
            const { data: userData } = await supabase
                .from("users")
                .select("credit_balance")
                .eq("id", session.user.id)
                .single();
            if (userData) credits = userData.credit_balance;

            fetchHistory(session.user.id);
            prompt = "";
        } catch (error: any) {
            console.error("Generation error:", error);
            alert(`Error: ${error.message}`);
        } finally {
            isGenerating = false;
        }
    }
</script>

<svelte:head>
    <title>Lesson Forge | AI Lesson Planner</title>
</svelte:head>

<Header
    title="Lesson Forge"
    {email}
    {credits}
    {isLoggedIn}
    onSignOut={handleSignOut}
/>

<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <!-- Hero Section -->
    <div class="text-center mb-16 animate-fade-in">
        <h2
            class="text-4xl md:text-5xl font-extrabold text-slate-900 tracking-tight mb-4"
        >
            Helping Forge Future <span class="text-primary italic">Minds</span>
        </h2>
        <p class="text-lg text-slate-500 max-w-2xl mx-auto">
            Generate high-quality presentations and lesson plans in seconds.
        </p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-10">
        <!-- Left: Input Form -->
        <div
            class="lg:col-span-8 space-y-8 animate-fade-in"
            style="animation-delay: 0.1s"
        >
            <div class="card-premium p-8 relative overflow-hidden">
                <div class="absolute top-0 right-0 p-4 opacity-5">
                    <svg
                        class="w-32 h-32 text-primary"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                </div>

                <div
                    class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8"
                >
                    <h3
                        class="text-xl font-bold text-slate-800 flex items-center gap-2"
                    >
                        Forge Magic
                    </h3>

                    <!-- Mode Switcher -->
                    <div class="flex p-1 bg-slate-100 rounded-xl relative z-10">
                        <button
                            class="px-4 py-2 text-sm font-bold rounded-lg transition-all {genMode ===
                            'lesson'
                                ? 'bg-white text-primary shadow-sm'
                                : 'text-slate-500 hover:text-slate-700'}"
                            onclick={() => (genMode = "lesson")}
                        >
                            Lesson Plan
                        </button>
                        <button
                            class="px-4 py-2 text-sm font-bold rounded-lg transition-all {genMode ===
                            'ppt'
                                ? 'bg-white text-primary shadow-sm'
                                : 'text-slate-500 hover:text-slate-700'}"
                            onclick={() => (genMode = "ppt")}
                        >
                            Presentation
                        </button>
                    </div>
                </div>

                <div class="space-y-6">
                    <div>
                        <label
                            for="prompt"
                            class="block text-sm font-bold text-slate-700 mb-2 uppercase tracking-wide"
                        >
                            {genMode === "lesson"
                                ? "Lesson Topic & Objectives"
                                : "Presentation Context & Requirements"}
                        </label>
                        <div class="relative group">
                            <textarea
                                id="prompt"
                                name="prompt"
                                rows="5"
                                bind:value={prompt}
                                class="w-full px-5 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:ring-4 focus:ring-primary/10 focus:border-primary focus:bg-white transition-all text-slate-800 placeholder-slate-400"
                                placeholder={genMode === "lesson"
                                    ? "E.g., A lesson for beginner adults about grocery shopping, covering food vocabulary and polite requests..."
                                    : "E.g., A presentation for business English students about remote team management..."}
                            ></textarea>
                        </div>
                    </div>

                    <!-- Image Generation Toggle (PPT Only) -->
                    {#if genMode === "ppt"}
                        <div
                            class="p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 flex items-center justify-between group hover:border-accent transition-colors animate-fade-in {!isLoggedIn
                                ? 'opacity-50 grayscale pointer-events-none'
                                : ''}"
                        >
                            <div class="flex items-center gap-4">
                                <div
                                    class="w-10 h-10 rounded-xl bg-white flex items-center justify-center text-accent shadow-sm group-hover:scale-110 transition-transform"
                                >
                                    <svg
                                        class="w-6 h-6"
                                        fill="none"
                                        stroke="currentColor"
                                        viewBox="0 0 24 24"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                                        />
                                    </svg>
                                </div>
                                <div>
                                    <h4
                                        class="text-sm font-bold text-slate-800"
                                    >
                                        Include AI Images
                                    </h4>
                                    <p class="text-xs text-slate-500">
                                        Add 4-6 high-quality images to your
                                        presentation
                                    </p>
                                </div>
                            </div>
                            <label
                                class="relative inline-flex items-center cursor-pointer"
                            >
                                <input
                                    type="checkbox"
                                    bind:checked={includeImages}
                                    class="sr-only peer"
                                    disabled={!isLoggedIn}
                                />
                                <div
                                    class="w-11 h-6 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"
                                ></div>
                            </label>
                        </div>
                    {/if}

                    <div
                        class="flex items-center justify-between pt-4 border-t border-slate-100"
                    >
                        <div
                            class="flex items-center gap-2 text-sm text-slate-500"
                        >
                            {#if !isLoggedIn}
                                <span
                                    class="text-amber-600 font-bold text-[10px] uppercase tracking-wider bg-amber-50 px-3 py-1 rounded-full animate-pulse border border-amber-100"
                                    >Sign in to start forging</span
                                >
                            {/if}
                        </div>
                        <Button
                            text={isGenerating
                                ? "Generating..."
                                : `Generate ${genMode === "lesson" ? "Lesson Plan" : "Magic"} (${creditCost} Credit${creditCost > 1 ? "s" : ""})`}
                            isLoading={isGenerating}
                            disabled={!isLoggedIn || !prompt}
                            onclick={handleGenerate}
                        />
                    </div>
                </div>
            </div>

            <!-- History Section -->
            <div class="card-premium overflow-hidden">
                <div
                    class="px-8 py-6 border-b border-slate-100 flex items-center justify-between bg-slate-50/50"
                >
                    <h3
                        class="text-lg font-bold text-slate-800 flex items-center gap-2"
                    >
                        <svg
                            class="w-5 h-5 text-secondary"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                        </svg>
                        Recent Activity
                    </h3>
                    <button
                        class="text-sm font-semibold text-primary hover:underline"
                        >View All</button
                    >
                </div>
                <div class="divide-y divide-slate-50">
                    {#if history.length === 0}
                        <div class="px-8 py-12">
                            <EmptyState
                                message="Your generated materials will appear here."
                            />
                        </div>
                    {:else}
                        {#each history as item}
                            <div
                                class="px-8 py-6 flex items-center justify-between hover:bg-slate-50/50 transition-colors group"
                            >
                                <div class="flex items-center gap-4">
                                    <div
                                        class="w-10 h-10 rounded-xl flex items-center justify-center {item.type ===
                                        'lesson'
                                            ? 'bg-emerald-50 text-emerald-600'
                                            : 'bg-blue-50 text-blue-600'}"
                                    >
                                        {#if item.type === "lesson"}
                                            <svg
                                                class="w-5 h-5"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
                                            >
                                                <path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                                                />
                                            </svg>
                                        {:else}
                                            <svg
                                                class="w-5 h-5"
                                                fill="none"
                                                stroke="currentColor"
                                                viewBox="0 0 24 24"
                                            >
                                                <path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"
                                                />
                                            </svg>
                                        {/if}
                                    </div>
                                    <div>
                                        <h4
                                            class="text-sm font-bold text-slate-800 line-clamp-1"
                                        >
                                            {item.title}
                                        </h4>
                                        <p
                                            class="text-xs text-slate-500 uppercase tracking-wider font-semibold mt-0.5"
                                        >
                                            {item.type} • {item.date}
                                        </p>
                                    </div>
                                </div>
                                <div class="flex items-center gap-3">
                                    {#if item.type === "lesson" && !item.hasPPT}
                                        <button
                                            class="px-3 py-1.5 bg-accent/20 text-primary text-[10px] font-black uppercase rounded-lg hover:bg-accent/40 shadow-sm transition-all"
                                            aria-label="Generate Presentation"
                                            onclick={() => {
                                                genMode = "ppt";
                                                prompt = `Create a presentation based on: ${item.title}`;
                                                item.hasPPT = true;
                                                window.scrollTo({
                                                    top: 0,
                                                    behavior: "smooth",
                                                });
                                            }}
                                        >
                                            Gen PPT
                                        </button>
                                    {:else if item.type === "ppt" && !item.hasLesson}
                                        <button
                                            class="px-3 py-1.5 bg-emerald-100 text-emerald-700 text-[10px] font-black uppercase rounded-lg hover:bg-emerald-200 shadow-sm transition-all"
                                            aria-label="Generate Lesson Plan"
                                            onclick={() => {
                                                genMode = "lesson";
                                                prompt = `Create a detailed lesson plan for: ${item.title}`;
                                                item.hasLesson = true;
                                                window.scrollTo({
                                                    top: 0,
                                                    behavior: "smooth",
                                                });
                                            }}
                                        >
                                            Gen Lesson
                                        </button>
                                    {/if}
                                    <button
                                        class="p-2 text-slate-400 hover:text-primary transition-colors"
                                        aria-label="Download"
                                    >
                                        <svg
                                            class="w-5 h-5"
                                            fill="none"
                                            stroke="currentColor"
                                            viewBox="0 0 24 24"
                                        >
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                                            />
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        {/each}
                    {/if}
                </div>
            </div>
        </div>

        <!-- Right: Credits & Sidebar -->
        <div
            class="lg:col-span-4 space-y-8 animate-fade-in"
            style="animation-delay: 0.2s"
        >
            <div
                class="card-premium p-8 bg-primary relative text-white overflow-hidden group border-none"
            >
                <div
                    class="absolute -bottom-10 -right-10 w-40 h-40 bg-white/10 rounded-full blur-3xl group-hover:bg-white/20 transition-all duration-700"
                ></div>

                <h3 class="text-xl font-bold mb-2">Fuel Your Forge</h3>
                <p class="text-white/70 text-sm mb-8 leading-relaxed">
                    Unlock unlimited creativity with credit bundles. No
                    subscriptions, just pay for what you use.
                </p>

                <div class="space-y-4 relative z-10">
                    <button
                        onclick={() =>
                            (window.location.href =
                                "https://buy.stripe.com/9B600lb2D6951Io1JsbjW03")}
                        class="w-full bg-white text-primary font-bold py-4 rounded-2xl hover:bg-surface transition-all shadow-xl hover:-translate-y-1"
                    >
                        Get 10 Credits <span class="text-primary/50 mx-2"
                            >|</span
                        > $9.99
                    </button>
                    <div
                        class="text-center group-hover:scale-[1.02] transition-transform duration-500"
                    >
                        <button
                            onclick={() =>
                                (window.location.href =
                                    "https://buy.stripe.com/9B64gBb2D695eva3RAbjW04")}
                            class="w-full bg-accent text-primary font-extrabold py-5 rounded-2xl hover:bg-accent/90 transition-all shadow-2xl relative"
                        >
                            <span
                                class="absolute -top-3 left-1/2 -translate-x-1/2 bg-secondary text-white text-[10px] px-3 py-1 rounded-full uppercase tracking-widest font-black shadow-lg"
                                >Popular</span
                            >
                            Get 50 Credits
                            <span class="text-primary/50 mx-2">|</span> $39.99
                        </button>
                    </div>
                </div>

                <div
                    class="mt-8 flex items-center gap-4 text-xs font-medium text-white/60"
                >
                    <div class="flex -space-x-2">
                        <div
                            class="w-6 h-6 rounded-full bg-white/20 border-2 border-primary"
                        ></div>
                        <div
                            class="w-6 h-6 rounded-full bg-white/30 border-2 border-primary"
                        ></div>
                        <div
                            class="w-6 h-6 rounded-full bg-white/40 border-2 border-primary"
                        ></div>
                    </div>
                    <span>Joined by 2,000+ Teachers</span>
                </div>
            </div>

            <div
                class="card-premium p-6 border-dashed border-2 flex flex-col items-center text-center"
            >
                <div
                    class="w-12 h-12 rounded-full bg-slate-50 flex items-center justify-center mb-4"
                >
                    <svg
                        class="w-6 h-6 text-slate-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                </div>
                <h4 class="font-bold text-slate-800 text-sm">Need Help?</h4>
                <p class="text-xs text-slate-400 mt-1 mb-4">
                    Our AI assistant is here to help with your lesson prompts.
                </p>
                <button class="text-primary text-xs font-bold hover:underline"
                    >Read Guide</button
                >
            </div>
        </div>
    </div>
</main>
