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
    
    // Form Fields
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
        if (session) {
            isLoggedIn = true;
            email = session.user.email || [cite_start]""; [cite: 107]
            await refreshCredits();
            await fetchHistory();
        } else {
            isLoggedIn = false;
            [cite_start]email = ""; [cite: 108]
            credits = 0;
            history = [];
            showPreview = false;
        }
    }

    async function handleSignOut() {
        [cite_start]await supabase.auth.signOut(); [cite: 109]
    }

    async function refreshCredits() {
        const { data: { session } } = await supabase.auth.getSession();
        [cite_start]if (!session) return; [cite: 111]
        const res = await fetch("/api/user/credits", {
            headers: { "Authorization": `Bearer ${session.access_token}` }
        });
        [cite_start]const data = await res.json(); [cite: 112]
        credits = data.credits;
    }

    async function fetchHistory() {
        const { data } = await supabase
            .from('generations')
            .select('*')
            .order('created_at', { ascending: false });
        [cite_start]if (data) history = data; [cite: 113]
    }

    async function handleGenerate() {
        [cite_start]if (!isLoggedIn || !canGenerate) return; [cite: 113]
        isGenerating = true;
        showPreview = false;
        
        const { data: { session } } = await supabase.auth.getSession();
        
        [cite_start]// Triggers the backend to generate AI content, upload to storage, and insert into DB [cite: 115]
        const res = await fetch("/api/generate", {
            method: "POST",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": `Bearer ${session?.access_token}` 
            },
            body: JSON.stringify({ 
                prompt, 
                grade, 
                duration, 
                [cite_start]mode: genMode === "lesson" ? "pdf" : "ppt", [cite: 116]
                teacher_name: teacherName,
                [cite_start]class_name: className [cite: 117]
            })
        });
        
        if (res.ok) {
            const data = await res.json();
            // Display markdown in the preview area
            [cite_start]generatedMarkdown = data.raw_content; [cite: 119]
            showPreview = true;
            
            // Immediately refresh local state so the new item appears in history
            await refreshCredits();
            [cite_start]await fetchHistory(); [cite: 120]
        } else {
            const error = await res.json();
            [cite_start]alert(error.error || "Generation failed"); [cite: 121]
        }
        isGenerating = false;
    }

    function printDoc() {
        [cite_start]window.print(); [cite: 123]
    }
</script>

<div class="min-h-screen bg-[#F8FAFC]">
    <div class="no-print">
        <Header 
            title="Vaelia Forge" 
            {email} 
            {credits} 
            {isLoggedIn} 
            onSignOut={handleSignOut} 
        />
    </div>

    <main class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
            
            <div class="lg:col-span-8 space-y-8">
                <div class="no-print bg-white p-8 rounded-3xl shadow-sm border border-slate-200 space-y-6">
                    <div class="flex items-center justify-between">
                        [cite_start]<h2 class="text-2xl font-bold text-slate-800">Forge New Content</h2> [cite: 125]
                        <div class="flex bg-slate-100 p-1 rounded-xl">
                            <button onclick={() => genMode = "lesson"} class="px-4 py-2 rounded-lg text-sm font-bold {genMode === 'lesson' ? 'bg-white shadow text-primary' : 'text-slate-500'}">Lesson Plan</button>
                            [cite_start]<button onclick={() => genMode = "ppt"} class="px-4 py-2 rounded-lg text-sm font-bold {genMode === 'ppt' ? 'bg-white shadow text-primary' : 'text-slate-500'}">Presentation</button> [cite: 126, 127]
                        </div>
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        [cite_start]<input bind:value={teacherName} placeholder="Teacher Name" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" /> [cite: 128]
                        <input bind:value={className} placeholder="Class/Subject" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" />
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        [cite_start]<input bind:value={grade} placeholder="Grade Level" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" /> [cite: 129]
                        <input bind:value={duration} placeholder="Duration" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" />
                    </div>
                    
                    [cite_start]<textarea bind:value={prompt} placeholder="What should we teach today?" class="w-full h-32 p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary"></textarea> [cite: 130]
                    
                    <div class="flex justify-between items-center">
                        <div class="flex flex-col">
                            [cite_start]<p class="text-sm text-slate-500 font-medium">Cost: <span class="text-primary font-bold">{creditCost} Credit</span></p> [cite: 131]
                            {#if isLoggedIn && credits < creditCost}
                                <p class="text-xs text-red-500 font-bold">Insufficient Credits</p>
                            {/if}
                        </div>
                        <Button 
                            onclick={handleGenerate} 
                            text={isLoggedIn ? "Generate Preview" : "Sign in to Generate"} 
                            isLoading={isGenerating} 
                            disabled={!isLoggedIn || [cite_start]!canGenerate || isGenerating} [cite: 133, 134, 135]
                        />
                    </div>
                </div>

                {#if showPreview}
                    [cite_start]<div id="printable-area" class="bg-white p-12 lg:p-16 shadow-2xl rounded-sm border border-slate-200 animate-in fade-in duration-700"> [cite: 136]
                        <div class="border-b-2 border-slate-900 pb-6 mb-10 flex justify-between items-end">
                            <div>
                                [cite_start]<h1 class="text-4xl font-serif font-bold text-slate-900 tracking-tight uppercase">Lesson Plan</h1> [cite: 137]
                                <p class="text-sm font-medium text-slate-500 mt-1 italic">Generated via Vaelia Forge</p>
                            </div>
                            <div class="text-right text-sm space-y-1 text-slate-700 font-mono uppercase">
                                <p><span class="font-bold">Teacher:</span> {teacherName || [cite_start]'____________'}</p> [cite: 138, 139]
                                <p><span class="font-bold">Class:</span> {className || [cite_start]'____________'}</p> [cite: 140]
                                <p><span class="font-bold">Grade:</span> {grade || [cite_start]'____________'}</p> [cite: 141]
                                <p><span class="font-bold">Date:</span> {new Date().toLocaleDateString()}</p>
                            </div>
                        </div>

                        <div class="prose prose-slate max-w-none">
                            [cite_start]{@html marked.parse(generatedMarkdown)} [cite: 142]
                        </div>
                    </div>
                    
                    <div class="flex justify-center no-print">
                        <button onclick={printDoc} class="bg-primary text-white px-10 py-5 rounded-2xl font-bold shadow-2xl hover:scale-105 active:scale-95 transition-all">
                            [cite_start]Download Stylized PDF [cite: 143, 144]
                        </button>
                    </div>
                {:else if !isGenerating}
                    <div class="no-print">
                        [cite_start]<EmptyState message="Your forged lesson will appear here..." /> [cite: 145]
                    </div>
                {/if}
            </div>

            <div class="lg:col-span-4 space-y-8 no-print">
                <div class="p-8 bg-primary rounded-3xl text-white shadow-xl">
                    [cite_start]<h3 class="text-xl font-bold mb-2">Get More Credits</h3> [cite: 146]
                    <p class="text-xs text-white/70 mb-6">Current Balance: {credits} Forges</p>
                    <div class="space-y-4 mt-6">
                        <button onclick={() => window.location.href = 'https://buy.stripe.com/9B600lb2D6951Io1JsbjW03'} class="w-full bg-white text-primary font-bold py-4 rounded-2xl shadow-md hover:-translate-y-1 transition-all">
                            10 Credits | [cite_start]$9.99 [cite: 147, 148]
                        </button>
                        <button onclick={() => window.location.href = 'https://buy.stripe.com/9B64gBb2D695eva3RAbjW04'} class="w-full bg-accent text-primary font-extrabold py-5 rounded-2xl shadow-lg relative hover:-translate-y-1 transition-all">
                            [cite_start]<span class="absolute -top-3 left-1/2 -translate-x-1/2 bg-white text-primary text-[10px] py-0.5 rounded-full border border-accent uppercase tracking-tighter">Best Value</span> [cite: 149]
                            25 Credits | [cite_start]$19.99 [cite: 150]
                        </button>
                    </div>
                </div>

                <div class="bg-white p-8 rounded-3xl shadow-sm border border-slate-200">
                    [cite_start]<h3 class="text-lg font-bold text-slate-800 mb-6">Past Forges</h3> [cite: 151]
                    <div class="space-y-4">
                        {#if history.length === 0}
                            <p class="text-slate-400 text-sm italic">No history yet.</p>
                        {:else}
                            {#each history as item}
                                <div class="p-4 bg-slate-50 rounded-2xl border border-slate-100 flex items-center justify-between group">
                                    [cite_start]<div class="truncate mr-4"> [cite: 153]
                                        <p class="font-bold text-slate-800 truncate text-sm">{item.prompt}</p>
                                        [cite_start]<p class="text-[10px] text-slate-400 uppercase tracking-widest">{new Date(item.created_at).toLocaleDateString()}</p> [cite: 154]
                                    </div>
                                    <a href={item.file_path} target="_blank" rel="noreferrer" class="p-2 bg-white rounded-xl shadow-sm opacity-0 group-hover:opacity-100 transition-opacity">
                                        [cite_start]<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor"> [cite: 155]
                                            [cite_start]<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /> [cite: 156]
                                        </svg>
                                    </a>
                                </div>
                            {/each}
                        {/if}
                    </div>
                </div>
            </div>
        </div>
    </main>
</div>

<style>
    @media print {
        :global(.no-print) { display: none !important; [cite_start]} [cite: 159]
        :global(body) { background: white !important; margin: 0; [cite_start]} [cite: 160]
        #printable-area { border: none !important; box-shadow: none !important; margin: 0 !important; padding: 0 !important; [cite_start]} [cite: 161]
    }
    
    :global(.prose h2) { 
        border-bottom: 1px solid #e2e8f0;
        [cite_start]padding-bottom: 0.5rem; [cite: 162]
        margin-top: 2rem;
        color: #0f172a;
        font-weight: 700;
    }
</style>