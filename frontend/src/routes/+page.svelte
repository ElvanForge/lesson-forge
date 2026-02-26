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

    let creditCost = $derived(genMode === "lesson" ? 1 : 1);

    onMount(() => {
        if (!isSupabaseConfigured) return;
        supabase.auth.getSession().then(({ data: { session } }) => handleAuthStateChange(session));
        const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => handleAuthStateChange(session));
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
        const response = await fetch("/api/user/credits", {
            headers: { "Authorization": `Bearer ${session.access_token}` }
        });
        if (response.ok) {
            const data = await response.json();
            credits = data.credits;
        }
    }

    async function fetchHistory() {
        const { data: { user } } = await supabase.auth.getUser();
        if (!user) return;
        const { data } = await supabase.from('generations').select('*').eq('user_id', user.id).order('created_at', { ascending: false });
        if (data) history = data;
    }

    async function handleGenerate() {
        if (!isLoggedIn || !prompt || isGenerating) return;
        if (credits < creditCost) {
            alert("Insufficient credits!");
            return;
        }
        
        isGenerating = true;
        downloadUrl = "";

        try {
            const { data: { session } } = await supabase.auth.getSession();
            const response = await fetch("/api/generate", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${session?.access_token}`
                },
                body: JSON.stringify({ prompt, mode: genMode === "lesson" ? "pdf" : "ppt" })
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || "Server error");
            }

            const data = await response.json();
            downloadUrl = data.file;
            await refreshCredits();
            await fetchHistory();
        } catch (err: any) {
            alert(`Forging failed: ${err.message}`);
        } finally {
            isGenerating = false;
        }
    }
</script>

<div class="min-h-screen bg-[#F8FAFC]">
    <Header {isLoggedIn} {credits} {email} title="Lesson Forge" onSignOut={() => supabase.auth.signOut()} />

    <main class="max-w-7xl mx-auto px-4 py-12">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12">
            <div class="lg:col-span-8 space-y-8">
                <div class="bg-white rounded-3xl p-8 shadow-sm border border-slate-100">
                    <h2 class="text-2xl font-bold text-slate-900 mb-6">Create Content</h2>
                    <div class="space-y-6">
                        <div class="flex p-1 bg-slate-100 rounded-2xl w-fit">
                            <button onclick={() => genMode = "lesson"} class="px-6 py-2 rounded-xl text-sm font-medium {genMode === 'lesson' ? 'bg-white text-primary shadow-sm' : 'text-slate-500'}">Lesson Plan</button>
                            <button onclick={() => genMode = "slides"} class="px-6 py-2 rounded-xl text-sm font-medium {genMode === 'slides' ? 'bg-white text-primary shadow-sm' : 'text-slate-500'}">Presentation</button>
                        </div>

                        <textarea bind:value={prompt} placeholder="Describe your lesson topic..." class="w-full h-40 p-6 bg-slate-50 border-none rounded-3xl resize-none"></textarea>

                        <div class="flex items-center justify-between bg-slate-50 p-4 rounded-2xl">
                            <div>
                                <p class="text-sm font-bold text-slate-900">Cost: {creditCost} Credits</p>
                                <p class="text-xs text-slate-500">Balance: {credits}</p>
                            </div>
                            <Button onclick={handleGenerate} disabled={isGenerating || !prompt} isLoading={isGenerating} text={isGenerating ? "Forging..." : "Forge Content"} />
                        </div>

                        {#if downloadUrl}
                            <div class="animate-bounce flex justify-center">
                                <a href={downloadUrl} target="_blank" class="bg-accent text-primary font-bold py-3 px-8 rounded-2xl shadow-lg">Download Ready</a>
                            </div>
                        {/if}
                    </div>
                </div>

                <div class="space-y-4">
                    <h3 class="text-xl font-bold text-slate-900">History</h3>
                    {#if history.length === 0}
                        <EmptyState message="Your forged files will appear here." />
                    {:else}
                        {#each history as item}
                            <div class="bg-white p-6 rounded-3xl border border-slate-100 flex items-center justify-between">
                                <p class="font-bold text-slate-900 truncate max-w-xs">{item.prompt}</p>
                                <a href={item.file_path} target="_blank" class="p-3 rounded-xl bg-slate-50 text-primary">Download</a>
                            </div>
                        {/each}
                    {/if}
                </div>
            </div>
        </div>
    </main>
</div>