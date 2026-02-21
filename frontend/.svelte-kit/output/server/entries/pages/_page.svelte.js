import { e as escape_html, $ as head, a0 as attr_class, a1 as attr, a2 as ensure_array_like, a3 as stringify, _ as derived } from "../../chunks/index.js";
import "clsx";
import { B as Button } from "../../chunks/Button.js";
import "../../chunks/supabase.js";
function Header($$renderer, $$props) {
  let { title, credits } = $$props;
  $$renderer.push(`<header class="glass sticky top-0 z-50 border-b border-white/40"><div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-20 flex items-center justify-between"><a href="/" class="flex items-center gap-3 hover:opacity-80 transition-opacity"><div class="w-10 h-10 bg-primary rounded-xl flex items-center justify-center shadow-lg shadow-primary/20"><svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg></div> <h1 class="text-2xl font-extrabold tracking-tight text-primary">${escape_html(title)}</h1></a> <div class="flex items-center gap-6">`);
  {
    $$renderer.push("<!--[!-->");
  }
  $$renderer.push(`<!--]--> <div class="flex items-center gap-3"><span class="flex items-center gap-1.5 text-sm font-bold bg-primary/10 text-primary px-4 py-2 rounded-full border border-primary/10"><span class="w-2 h-2 rounded-full bg-primary animate-pulse"></span> ${escape_html(credits)} Credits</span> `);
  {
    $$renderer.push("<!--[!-->");
    $$renderer.push(`<a href="/login" class="text-sm font-bold text-primary hover:text-primary/80 transition-all">Sign In</a>`);
  }
  $$renderer.push(`<!--]--></div></div></div></header>`);
}
function EmptyState($$renderer, $$props) {
  let { message } = $$props;
  $$renderer.push(`<div class="text-center py-20 animate-fade-in"><div class="inline-flex items-center justify-center w-20 h-20 rounded-3xl bg-slate-50 mb-6 text-slate-300"><svg class="h-10 w-10" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg></div> <h3 class="text-lg font-bold text-slate-800">No Generations Yet</h3> <p class="mt-2 text-slate-500 max-w-xs mx-auto">${escape_html(message)}</p></div>`);
}
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let credits = 0;
    let isGenerating = false;
    let prompt = "";
    let includeImages = false;
    let history = [];
    let creditCost = derived(() => 1);
    async function handleGenerate() {
      return;
    }
    head("1uha8ag", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>Lesson Forge | AI Lesson Planner</title>`);
      });
    });
    Header($$renderer2, {
      title: "Lesson Forge",
      credits
    });
    $$renderer2.push(`<!----> <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12"><div class="text-center mb-16"><h2 class="text-4xl md:text-5xl font-extrabold text-slate-900 tracking-tight mb-4">Helping Forge Future <span class="text-primary italic">Minds</span></h2> <p class="text-lg text-slate-500 max-w-2xl mx-auto">Generate high-quality presentations and lesson plans in seconds.</p></div> <div class="grid grid-cols-1 lg:grid-cols-12 gap-10"><div class="lg:col-span-8 space-y-8"><div class="card-premium p-8 relative overflow-hidden bg-white border rounded-3xl shadow-sm"><div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8"><h3 class="text-xl font-bold text-slate-800 flex items-center gap-2">Forge Magic</h3> <div class="flex p-1 bg-slate-100 rounded-xl relative z-10"><button${attr_class(`px-4 py-2 text-sm font-bold rounded-lg transition-all ${stringify("text-slate-500")}`)}>Lesson Plan</button> <button${attr_class(`px-4 py-2 text-sm font-bold rounded-lg transition-all ${stringify("bg-white text-primary shadow-sm")}`)}>Presentation</button></div></div> <div class="space-y-6"><div><label for="prompt" class="block text-sm font-bold text-slate-700 mb-2 uppercase tracking-wide">${escape_html("Presentation Context")}</label> <textarea id="prompt" rows="5" class="w-full px-5 py-4 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:ring-4 focus:ring-primary/10 focus:border-primary focus:bg-white transition-all text-slate-800"${attr("placeholder", "E.g., Presentation about remote team management...")}>`);
    const $$body = escape_html(prompt);
    if ($$body) {
      $$renderer2.push(`${$$body}`);
    }
    $$renderer2.push(`</textarea></div> `);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div class="p-4 bg-slate-50 rounded-2xl border-2 border-slate-100 flex items-center justify-between group hover:border-accent transition-colors"><div class="flex items-center gap-4"><div class="w-10 h-10 rounded-xl bg-white flex items-center justify-center text-accent shadow-sm"><svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg></div> <div><h4 class="text-sm font-bold text-slate-800">Include AI Images</h4> <p class="text-xs text-slate-500">Add 4-6 high-quality images</p></div></div> <label class="relative inline-flex items-center cursor-pointer"><input type="checkbox"${attr("checked", includeImages, true)} class="sr-only peer"/> <div class="w-11 h-6 bg-slate-200 rounded-full peer peer-checked:after:translate-x-full peer-checked:bg-primary after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all"></div></label></div>`);
    }
    $$renderer2.push(`<!--]--> <div class="flex items-center justify-between pt-4 border-t border-slate-100"><div class="text-sm text-slate-500">`);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<span class="text-amber-600 font-bold text-[10px] uppercase tracking-wider bg-amber-50 px-3 py-1 rounded-full">Sign in to start forging</span>`);
    }
    $$renderer2.push(`<!--]--></div> `);
    Button($$renderer2, {
      text: `Generate (${creditCost()} Credit${creditCost() > 1 ? "s" : ""})`,
      isLoading: isGenerating,
      disabled: true,
      onclick: handleGenerate
    });
    $$renderer2.push(`<!----></div></div></div> <div class="card-premium overflow-hidden bg-white border rounded-3xl shadow-sm"><div class="px-8 py-6 border-b border-slate-100 flex items-center justify-between bg-slate-50/50"><h3 class="text-lg font-bold text-slate-800 flex items-center gap-2">Recent Activity</h3></div> <div class="divide-y divide-slate-50">`);
    if (history.length === 0) {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div class="px-8 py-12">`);
      EmptyState($$renderer2, { message: "Your forged materials will appear here." });
      $$renderer2.push(`<!----></div>`);
    } else {
      $$renderer2.push("<!--[!-->");
      $$renderer2.push(`<!--[-->`);
      const each_array = ensure_array_like(history);
      for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
        let item = each_array[$$index];
        $$renderer2.push(`<div class="px-8 py-6 flex items-center justify-between hover:bg-slate-50/50 transition-colors group"><div class="flex items-center gap-4"><div${attr_class(`w-10 h-10 rounded-xl flex items-center justify-center ${stringify(item.type === "lesson" ? "bg-emerald-50 text-emerald-600" : "bg-blue-50 text-blue-600")}`)}>`);
        if (item.type === "lesson") {
          $$renderer2.push("<!--[-->");
          $$renderer2.push(`<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>`);
        } else {
          $$renderer2.push("<!--[!-->");
          $$renderer2.push(`<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"></path></svg>`);
        }
        $$renderer2.push(`<!--]--></div> <div><h4 class="text-sm font-bold text-slate-800 line-clamp-1">${escape_html(item.title)}</h4> <p class="text-xs text-slate-500 uppercase font-semibold">${escape_html(item.type)} • ${escape_html(item.date)}</p></div></div> <div class="flex items-center gap-3"><button class="p-2 text-slate-400 hover:text-primary transition-colors"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a2 2 0 002 2h12a2 2 0 002-2v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg></button></div></div>`);
      }
      $$renderer2.push(`<!--]-->`);
    }
    $$renderer2.push(`<!--]--></div></div></div> <div class="lg:col-span-4 space-y-8"><div class="p-8 bg-primary rounded-3xl text-white shadow-xl"><h3 class="text-xl font-bold mb-2">Fuel Your Forge</h3> <p class="text-white/70 text-sm mb-8">Unlock unlimited creativity with credit bundles.</p> <div class="space-y-4"><button class="w-full bg-white text-primary font-bold py-4 rounded-2xl shadow-xl hover:-translate-y-1 transition-all">10 Credits | $9.99</button> <button class="w-full bg-accent text-primary font-extrabold py-5 rounded-2xl shadow-2xl relative"><span class="absolute -top-3 left-1/2 -translate-x-1/2 bg-secondary text-white text-[10px] px-3 py-1 rounded-full">Popular</span> 50 Credits | $39.99</button></div></div></div></div></main>`);
  });
}
export {
  _page as default
};
