<script lang="ts">
    import { marked } from 'marked';
    import Header from "$lib/components/Header.svelte";
    import Button from "$lib/components/Button.svelte";
    import EmptyState from "$lib/components/EmptyState.svelte";
    import { onMount } from "svelte";
    import { supabase, isSupabaseConfigured } from "$lib/supabase";

    [cite_start]let isLoggedIn = $state(false); [cite: 74]
    [cite_start]let credits = $state(0); [cite: 74]
    [cite_start]let email = $state(""); [cite: 74]
    [cite_start]let isGenerating = $state(false); [cite: 74]
    
    [cite_start]let prompt = $state(""); [cite: 74]
    [cite_start]let grade = $state(""); [cite: 75]
    [cite_start]let duration = $state(""); [cite: 75]
    [cite_start]let teacherName = $state(""); [cite: 75]
    [cite_start]let className = $state(""); [cite: 75]
    
    [cite_start]let genMode = $state("lesson"); [cite: 75]
    [cite_start]let history = $state<any[]>([]); [cite: 76]
    [cite_start]let generatedMarkdown = $state(""); [cite: 76]
    [cite_start]let showPreview = $state(false); [cite: 76]

    [cite_start]let creditCost = $derived(genMode === "lesson" ? 1 : 2); [cite: 77]
    [cite_start]let canGenerate = $derived(credits >= creditCost && prompt.length > 0); [cite: 78]

    onMount(() => {
        [cite_start]if (!isSupabaseConfigured) return; [cite: 79]
        supabase.auth.getSession().then(({ data: { session } }) => {
            [cite_start]handleAuthStateChange(session); [cite: 79]
        });
        const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => {
            [cite_start]handleAuthStateChange(session); [cite: 79]
        });
        [cite_start]return () => subscription.unsubscribe(); [cite: 79]
    });

    async function handleAuthStateChange(session: any) {
        [cite_start]if (session && session.user) { [cite: 80]
            [cite_start]isLoggedIn = true; [cite: 80]
            email = session.user.email ?? [cite_start]""; [cite: 81]
            [cite_start]await refreshCredits(); [cite: 81]
            [cite_start]await fetchHistory(); [cite: 81]
        } else {
            [cite_start]isLoggedIn = false; [cite: 81]
            [cite_start]email = ""; [cite: 82]
            [cite_start]credits = 0; [cite: 82]
            [cite_start]history = []; [cite: 82]
            [cite_start]showPreview = false; [cite: 82]
        }
    }

    async function handleSignOut() {
        [cite_start]await supabase.auth.signOut(); [cite: 83]
    }

    async function refreshCredits() {
        [cite_start]const { data: { session } } = await supabase.auth.getSession(); [cite: 84]
        [cite_start]if (!session) return; [cite: 85]
        const res = await fetch("/api/user/credits", {
            [cite_start]headers: { "Authorization": `Bearer ${session.access_token}` } [cite: 85]
        });
        [cite_start]const data = await res.json(); [cite: 86]
        [cite_start]credits = data.credits; [cite: 86]
    }

    async function fetchHistory() {
        [cite_start]const { data } = await supabase [cite: 86]
            [cite_start].from('generations') [cite: 87]
            [cite_start].select('*') [cite: 87]
            [cite_start].order('created_at', { ascending: false }); [cite: 87]
        [cite_start]if (data) history = data; [cite: 87]
    }

    async function handleGenerate() {
        [cite_start]if (!isLoggedIn || !canGenerate) return; [cite: 87]
        [cite_start]isGenerating = true; [cite: 88]
        [cite_start]showPreview = false; [cite: 88]
        
        [cite_start]const { data: { session } } = await supabase.auth.getSession(); [cite: 88]
        [cite_start]const res = await fetch("/api/generate", { [cite: 89]
            method: "POST",
            headers: { 
                "Content-Type": "application/json",
                [cite_start]"Authorization": `Bearer ${session?.access_token}` [cite: 89]
            },
            body: JSON.stringify({ 
                [cite_start]prompt, grade, duration, mode: genMode, [cite: 90]
                [cite_start]teacher_name: teacherName, class_name: className [cite: 90]
            })
        });

        [cite_start]if (res.ok) { [cite: 91]
            [cite_start]const data = await res.json(); [cite: 91]
            [cite_start]generatedMarkdown = data.raw_content; [cite: 92]
            [cite_start]showPreview = true; [cite: 92]
            [cite_start]await refreshCredits(); [cite: 92]
            [cite_start]await fetchHistory(); [cite: 92]
        }
        [cite_start]isGenerating = false; [cite: 93]
    }

    function printDoc() {
        [cite_start]window.print(); [cite: 94]
    }
</script>

<style>
    /* Fixed Styling for Stylized PDFs */
    @media print {
        /* Hide all UI elements except the content */
        :global(.no-print), :global(header), :global(.lg\:col-span-4) {
            display: none !important;
        }
        :global(body) {
            background: white !important;
        }
        .printable-content {
            box-shadow: none !important;
            border: none !important;
            padding: 0 !important;
            margin: 0 !important;
            width: 100% !important;
        }
    }
</style>

<div class="min-h-screen bg-[#F8FAFC]">
    <div class="no-print">
        <Header title="Vaelia Forge" {email} {credits} {isLoggedIn} onSignOut={handleSignOut} />
    </div>

    <main class="max-w-7xl mx-auto py-12 px-4">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
            <div class="lg:col-span-8 space-y-8">
                
                <div class="no-print bg-white p-8 rounded-3xl shadow-sm border border-slate-200 space-y-6">
                    <div class="flex items-center justify-between">
                        [cite_start]<h2 class="text-2xl font-bold text-slate-800">Forge New Content</h2> [cite: 96]
                        <div class="flex bg-slate-100 p-1 rounded-xl">
                            [cite_start]<button onclick={() => genMode = "lesson"} class="px-4 py-2 rounded-lg text-sm font-bold {genMode === 'lesson' ? 'bg-white shadow text-primary' : 'text-slate-500'}">Lesson Plan</button> [cite: 96, 97]
                            [cite_start]<button onclick={() => genMode = "ppt"} class="px-4 py-2 rounded-lg text-sm font-bold {genMode === 'ppt' ? 'bg-white shadow text-primary' : 'text-slate-500'}">Presentation</button> [cite: 97, 98]
                        </div>
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        [cite_start]<input bind:value={teacherName} placeholder="Teacher Name" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" /> [cite: 98, 99]
                        [cite_start]<input bind:value={className} placeholder="Class/Subject" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" /> [cite: 99]
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        [cite_start]<input bind:value={grade} placeholder="Grade Level" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" /> [cite: 100]
                        [cite_start]<input bind:value={duration} placeholder="Duration" class="p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary" /> [cite: 100]
                    </div>
                    
                    [cite_start]<textarea bind:value={prompt} placeholder="What should we teach today?" class="w-full h-32 p-4 bg-slate-50 rounded-2xl border-none focus:ring-2 ring-primary"></textarea> [cite: 101]
                    
                    <div class="flex justify-between items-center">
                        <div class="flex flex-col">
                            [cite_start]<p class="text-sm text-slate-500 font-medium">Cost: <span class="text-primary font-bold">{creditCost} Credit</span></p> [cite: 102]
                            {#if isLoggedIn && credits < creditCost}
                                [cite_start]<p class="text-xs text-red-500 font-bold">Insufficient Credits</p> [cite: 102]
                            {/if}
                        </div>
                        <Button onclick={handleGenerate} text={isLoggedIn ? "Generate Preview" : "Sign in to Generate"} isLoading={isGenerating} disabled={!isLoggedIn || [cite_start]!canGenerate || isGenerating} /> [cite: 103, 104, 105]
                    </div>
                </div>

                {#if isLoggedIn && credits < 5}
                    <div class="no-print p-6 bg-gradient-to-br from-indigo-50 to-white border border-indigo-100 rounded-3xl shadow-sm flex items-center justify-between">
                        <div class="flex items-center gap-4">
                            <div class="bg-indigo-500 p-3 rounded-2xl text-white">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                                </svg>
                            </div>
                            <div>
                                <h4 class="font-bold text-slate-800">Add More Forge Power</h4>
                                <p class="text-sm text-slate-500">Upgrade your account to keep creating high-quality lessons.</p>
                            </div>
                        </div>
                        <a href="https://buy.stripe.com/your_actual_link" target="_blank" class="bg-indigo-600 text-white px-6 py-3 rounded-xl font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-200">
                            Upgrade Now
                        </a>
                    </div>
                {/if}

                {#if showPreview}
                    <div class="printable-content bg-white p-12 lg:p-16 shadow-2xl rounded-sm border border-slate-200">
                        [cite_start]<div class="prose prose-slate max-w-none"> [cite: 106]
                            <div class="border-b-4 border-primary pb-4 mb-8">
                                <h1 class="text-4xl font-serif font-bold text-slate-900 tracking-tight uppercase">
                                    {genMode === 'ppt' ? [cite_start]'Presentation Preview' : 'Lesson Plan'} [cite: 106, 107]
                                </h1>
                            </div>
                            [cite_start]{@html marked.parse(generatedMarkdown.replace(/---/g, '<hr class="my-8 border-slate-200" />'))} [cite: 107]
                        </div>
                    </div>
                    [cite_start]<div class="flex justify-center no-print mt-8"> [cite: 108]
                        {#if genMode === 'ppt'}
                             [cite_start]<a href={history[0].file_path} download class="bg-primary text-white px-10 py-5 rounded-2xl font-bold shadow-2xl">Download Presentation (.pptx)</a> [cite: 108]
                        {:else}
                            [cite_start]<button onclick={printDoc} class="bg-primary text-white px-10 py-5 rounded-2xl font-bold shadow-2xl">Print Stylized PDF</button> [cite: 109]
                        {/if}
                    </div>
                {:else if !isGenerating}
                    [cite_start]<EmptyState message="Your forged content will appear here..." /> [cite: 110]
                {/if}
            </div>

            [cite_start]<div class="lg:col-span-4 space-y-8 no-print"> [cite: 110]
                [cite_start]<div class="bg-white p-8 rounded-3xl shadow-sm border border-slate-200"> [cite: 111]
                    [cite_start]<h3 class="text-lg font-bold text-slate-800 mb-6">Past Forges</h3> [cite: 111]
                    [cite_start]<div class="space-y-4"> [cite: 111]
                        [cite_start]{#each history as item} [cite: 111]
                            [cite_start]<div class="p-4 bg-slate-50 rounded-2xl border border-slate-100 flex items-center justify-between group"> [cite: 112]
                                [cite_start]<div class="truncate mr-4"> [cite: 112]
                                    [cite_start]<p class="font-bold text-slate-800 truncate text-sm">{item.prompt}</p> [cite: 112]
                                    [cite_start]<p class="text-[10px] text-slate-400 uppercase tracking-widest">{new Date(item.created_at).toLocaleDateString()}</p> [cite: 113]
                                </div>
                                [cite_start]<a href={item.file_path} download target="_blank" rel="noreferrer" class="p-2 bg-white rounded-xl shadow-sm opacity-0 group-hover:opacity-100 transition-opacity"> [cite: 114]
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        [cite_start]<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /> [cite: 115]
                                    </svg>
                                </a>
                            [cite_start]</div> [cite: 116]
                        {/each}
                    </div>
                </div>
            </div>
        </div>
    </main>
</div>