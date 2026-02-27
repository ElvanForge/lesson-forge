<script lang="ts">
    import { marked } from 'marked';
    import Header from "$lib/components/Header.svelte";
    import Button from "$lib/components/Button.svelte";
    import EmptyState from "$lib/components/EmptyState.svelte";
    import { onMount } from "svelte";
    import { supabase, isSupabaseConfigured } from "$lib/supabase";

    let isLoggedIn = $state(false);
    let credits = $state(0);
    let isGenerating = $state(false);
    let prompt = $state("");
    let grade = $state("");
    let duration = $state("");
    let genMode = $state("lesson");
    let generatedMarkdown = $state("");
    let showPreview = $state(false);
    let history = $state<any[]>([]);

    onMount(async () => {
        const { data: { session } } = await supabase.auth.getSession();
        if (session) {
            isLoggedIn = true;
            refreshCredits();
            fetchHistory();
        }
    });

    async function refreshCredits() {
        const { data: { session } } = await supabase.auth.getSession();
        const res = await fetch("/api/user/credits", {
            headers: { "Authorization": `Bearer ${session?.access_token}` }
        });
        const data = await res.json();
        credits = data.credits;
    }

    async function fetchHistory() {
        const { data } = await supabase.from('generations').select('*').order('created_at', { ascending: false });
        if (data) history = data;
    }

    async function handleGenerate() {
        isGenerating = true;
        showPreview = false;
        const { data: { session } } = await supabase.auth.getSession();
        const res = await fetch("/api/generate", {
            method: "POST",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": `Bearer ${session?.access_token}` 
            },
            body: JSON.stringify({ prompt, grade, duration, mode: genMode === "lesson" ? "pdf" : "ppt" })
        });
        const data = await res.json();
        generatedMarkdown = data.raw_content;
        showPreview = true;
        isGenerating = false;
        refreshCredits();
        fetchHistory();
    }

    function printDoc() { window.print(); }
</script>

<div class="min-h-screen bg-[#F8FAFC]">
    <div class="no-print">
        <Header {isLoggedIn} {credits} title="Vaelia Forge" />
    </div>

    <main class="max-w-5xl mx-auto py-12 px-4">
        <div class="no-print bg-white p-8 rounded-3xl shadow-sm border mb-12 space-y-6">
            <h2 class="text-2xl font-bold">Forge Lesson Plan</h2>
            <div class="grid grid-cols-2 gap-4">
                <input bind:value={grade} placeholder="Grade Level" class="p-4 bg-slate-50 rounded-2xl border-none" />
                <input bind:value={duration} placeholder="Duration" class="p-4 bg-slate-50 rounded-2xl border-none" />
            </div>
            <textarea bind:value={prompt} placeholder="Enter your topic..." class="w-full h-32 p-4 bg-slate-50 rounded-2xl border-none"></textarea>
            <Button onclick={handleGenerate} text="Generate Preview" isLoading={isGenerating} />
        </div>

        {#if showPreview}
            <div id="printable-area" class="bg-white p-12 shadow-2xl rounded-sm border border-slate-200">
                <div class="text-center border-b-4 border-slate-900 pb-6 mb-8">
                    <h1 class="text-4xl font-serif font-bold text-slate-900 tracking-tighter">VAELIA FORGE</h1>
                    <p class="text-xs font-black tracking-[0.3em] uppercase mt-2 text-slate-500">Premium Education Resource</p>
                </div>

                <div class="prose prose-slate max-w-none prose-h1:text-3xl prose-h2:text-primary prose-h2:border-l-4 prose-h2:border-primary prose-h2:pl-4 prose-hr:border-slate-100">
                    {@html marked(generatedMarkdown)}
                </div>
            </div>
            
            <div class="flex justify-center mt-8 no-print">
                <button onclick={printDoc} class="bg-primary text-white px-8 py-4 rounded-2xl font-bold shadow-lg hover:scale-105 transition-transform">
                    Download Stylized PDF
                </button>
            </div>
        {/if}
    </main>
</div>

<style>
    @media print {
        .no-print { display: none !important; }
        body { background: white !important; }
        #printable-area { border: none !important; shadow: none !important; margin: 0 !important; padding: 0 !important; }
    }
    :global(.prose h2) { margin-top: 2.5rem; margin-bottom: 1rem; color: #1e293b; }
</style>