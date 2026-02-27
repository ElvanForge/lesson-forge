<script lang="ts">
    import { marked } from 'marked';
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
    let grade = $state("");
    let duration = $state("");
    let teacherName = $state("");
    let className = $state("");
    
    let genMode = $state("lesson");
    let history = $state<any[]>([]);
    let generatedMarkdown = $state("");
    let showPreview = $state(false);

    let creditCost = $derived(genMode === "lesson" ? 1 : 2);
    let canGenerate = $derived(credits >= creditCost && prompt.length > 0);

    onMount(() => {
        if (!isSupabaseConfigured) return;
        supabase.auth.getSession().then(({ data: { session } }) => {
            handleAuthStateChange(session);
        });
        const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => {
            handleAuthStateChange(session);
        });
        return () => subscription.unsubscribe();
    });

    async function handleAuthStateChange(session: any) {
        if (session && session.user) {
            isLoggedIn = true;
            email = session.user.email ?? ""; 
            await refreshCredits();
            await fetchHistory();
        } else {
            isLoggedIn = false;
            email = "";
            credits = 0;
            history = [];
            showPreview = false;
        }
    }

    async function handleSignOut() {
        await supabase.auth.signOut();
    }

    async function refreshCredits() {
        const { data: { session } } = await supabase.auth.getSession();
        if (!session) return;
        const res = await fetch("/api/user/credits", {
            headers: { "Authorization": `Bearer ${session.access_token}` }
        });
        const data = await res.json();
        credits = data.credits;
    }

    async function fetchHistory() {
        const { data } = await supabase
            .from('generations')
            .select('*')
            .order('created_at', { ascending: false });
        if (data) history = data;
    }

    async function handleGenerate() {
        if (!isLoggedIn || !canGenerate) return;
        isGenerating = true;
        showPreview = false;
        
        const { data: { session } } = await supabase.auth.getSession();
        const res = await fetch("/api/generate", {
            method: "POST",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": `Bearer ${session?.access_token}` 
            },
            body: JSON.stringify({ 
                prompt, grade, duration, mode: genMode,
                teacher_name: teacherName, class_name: className
            })
        });
        
        if (res.ok) {
            const data = await res.json();
            generatedMarkdown = data.raw_content;
            showPreview = true;
            await refreshCredits();
            await fetchHistory();
        }
        isGenerating = false;
    }

    function printDoc() {
        window.print();
    }
</script>

<div class="min-h-screen bg-[#F8FAFC]">
    <div class="no-print">
        <Header title="Vaelia Forge" {email} {credits} {isLoggedIn} onSignOut={handleSignOut} />
    </div>

    <main class="max-w-7xl mx-auto py-12 px-4">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
            <div class="lg:col-span-8 space-y-8">
                <div class="no-print bg-white p-8 rounded-3xl shadow-sm border border-slate-200 space-y-6">
                    <div class="flex items-center justify-between">
                        <h2 class="text-2xl font-bold text-slate-800">Forge New Content</h2>
                        <div class="flex bg-slate-100 p-1 rounded-xl">
                            <button onclick={() => genMode = "lesson"} class="px-4 py-2 rounded-lg text-sm font-bold {genMode === 'lesson' ? 'bg-white shadow text-primary' : 'text-slate-500'}">Lesson Plan</button>
                            <button onclick={() => genMode = "ppt"} class="px-4 py-2 rounded-lg text-sm font-bold {genMode === 'ppt' ? 'bg-white shadow text-primary' : 'text-slate-500'}">Presentation</button>
                        </div>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <input bind:value={teacherName} placeholder="Teacher Name" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" />
                        <input bind:value={className} placeholder="Class/Subject" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" />
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <input bind:value={grade} placeholder="Grade Level" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" />
                        <input bind:value={duration} placeholder="Duration" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" />
                    </div>
                    <textarea bind:value={prompt} placeholder="What should we teach today?" class="w-full h-32 p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary"></textarea>
                    
                    <div class="flex justify-between items-center">
                        <p class="text-sm text-slate-500 font-medium">Cost: <span class="text-primary font-bold">{creditCost} Credit</span></p>
                        <Button onclick={handleGenerate} text={isLoggedIn ? "Generate Preview" : "Sign in to Generate"} isLoading={isGenerating} disabled={!isLoggedIn || !canGenerate || isGenerating} />
                    </div>
                </div>

                {#if showPreview}
                    <div id="printable-area" class="bg-white p-12 lg:p-16 shadow-2xl rounded-sm border border-slate-200">
                        <div class="prose prose-slate max-w-none">
                            {@html marked.parse(generatedMarkdown)}
                        </div>
                    </div>
                    <div class="flex justify-center no-print mt-8">
                        <button onclick={printDoc} class="bg-primary text-white px-10 py-5 rounded-2xl font-bold">Download Stylized PDF</button>
                    </div>
                {:else if !isGenerating}
                    <EmptyState message="Your forged content will appear here..." />
                {/if}
            </div>

            <div class="lg:col-span-4 space-y-8 no-print">
                <div class="bg-white p-8 rounded-3xl shadow-sm border border-slate-200">
                    <h3 class="text-lg font-bold text-slate-800 mb-6">Past Forges</h3>
                    <div class="space-y-4">
                        {#each history as item}
                            <div class="p-4 bg-slate-50 rounded-2xl border border-slate-100 flex items-center justify-between group">
                                <div class="truncate mr-4">
                                    <p class="font-bold text-slate-800 truncate text-sm">{item.prompt}</p>
                                    <p class="text-[10px] text-slate-400 uppercase tracking-widest">{new Date(item.created_at).toLocaleDateString()}</p>
                                </div>
                                <a href={item.file_path} download target="_blank" rel="noreferrer" class="p-2 bg-white rounded-xl shadow-sm opacity-0 group-hover:opacity-100 transition-opacity">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                                    </svg>
                                </a>
                            </div>
                        {/each}
                    </div>
                </div>
            </div>
        </div>
    </main>
</div>

<style>
    @media print {
        :global(.no-print) { display: none !important; }
        :global(body) { background: white !important; margin: 0; }
        #printable-area { border: none !important; box-shadow: none !important; margin: 0 !important; padding: 0 !important; }
    }
</style>