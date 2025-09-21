<script lang="ts">
  import PageLayout from '../../components/PageLayout.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { onMount } from 'svelte';
  let htmlContent = $state('');

  onMount(() => {
    window.scrollTo(0, 0);
    fetch('/content/open-source.html')
      .then((res) => res.text())
      .then((data) => {
        htmlContent = data;
      });
  });
</script>

<svelte:head>
  <title>Open Source | {$LL.appName()}</title>
</svelte:head>

<div class="relative overflow-hidden bg-gradient-to-br from-blue-600 via-blue-700 to-indigo-900 dark:from-blue-900 dark:via-slate-900 dark:to-gray-900 flex items-center">
    <div class="absolute inset-0">
        <!-- Gradient Mesh -->
        <div class="absolute inset-0 bg-gradient-to-tr from-purple-500/20 via-transparent to-cyan-500/20"></div>  
        <!-- Radial Gradient Overlay -->
        <div class="absolute inset-0 bg-gradient-radial from-transparent via-transparent to-black/20"></div>
    </div>
    
    <!-- Content -->
    <div class="relative max-w-7xl mx-auto px-6 py-20 z-10">
        <div class="text-center text-white">            
            <!-- Main Heading -->
            <h1 class="text-6xl md:text-8xl font-black mb-8 leading-none tracking-tight">
                <span class="block transform transition-all duration-700 hover:scale-105">
                    Open Source
                </span>
                <span class="block bg-gradient-to-r from-blue-200 via-purple-200 to-cyan-200 bg-clip-text text-transparent relative transform transition-all duration-700 hover:scale-105">
                    First
                    <!-- Subtle glow effect -->
                    <div class="absolute inset-0 bg-gradient-to-r from-blue-400/20 via-purple-400/20 to-cyan-400/20 blur-3xl -z-10 animate-pulse"></div>
                </span>
            </h1>
            
            <!-- Subtitle -->
            <p class="text-2xl md:text-3xl font-light mb-8 opacity-90 max-w-4xl mx-auto transform transition-all duration-500 hover:opacity-100">
                Building <span class="text-blue-200 font-medium bg-gradient-to-r from-blue-200 to-cyan-200 bg-clip-text text-transparent">Agile Tools</span> the Right Way
            </p>
            
            <!-- Description -->
            <p class="text-lg md:text-xl opacity-80 max-w-5xl mx-auto leading-relaxed mb-12 transform transition-all duration-500 hover:opacity-95">
                At Thunderdome, open source isn't just a development model—it's a <strong class="text-white font-semibold">core principle</strong> that shapes everything we do. We've remained committed to building agile planning tools that are transparent, community-driven, and accessible to teams worldwide.
            </p>
            
            <!-- Enhanced CTA Buttons -->
            <div class="flex flex-col sm:flex-row gap-6 justify-center">
                <a href="https://thunderdome.dev" target="_blank" class="group relative inline-flex items-center gap-3 bg-white text-blue-600 px-10 py-5 rounded-2xl font-bold text-lg hover:bg-blue-50 transition-all duration-500 shadow-2xl shadow-blue-500/25 hover:shadow-blue-500/40 hover:scale-105 transform hover:-translate-y-1">
                    <span>Try Thunderdome</span>
                    <svg class="w-5 h-5 group-hover:translate-x-1 transition-transform duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                    </svg>
                    <!-- Button glow -->
                    <div class="absolute inset-0 bg-white rounded-2xl blur opacity-20 group-hover:opacity-30 transition-opacity duration-500 -z-10"></div>
                </a>
                
                <a href="https://github.com/StevenWeathers/thunderdome-planning-poker/blob/main/docs/INSTALLATION.md" target="_blank" class="group relative inline-flex items-center gap-3 bg-white/10 backdrop-blur-md text-white px-10 py-5 rounded-2xl font-bold text-lg hover:bg-white/20 transition-all duration-500 border border-white/30 hover:border-white/50 hover:scale-105 transform hover:-translate-y-1 overflow-hidden">
                    <svg class="w-5 h-5 group-hover:rotate-12 transition-transform duration-300" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                    <span>View on GitHub</span>
                    <!-- Shimmer effect using gradient -->
                    <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/10 to-transparent -skew-x-12 translate-x-[-200%] group-hover:translate-x-[200%] transition-transform duration-1000 ease-out"></div>
                </a>
            </div>
            
            <!-- Stats or Social Proof -->
            <div class="mt-16 flex flex-wrap justify-center gap-8 text-sm opacity-60 transition-opacity duration-300 hover:opacity-80">
                <div class="flex items-center gap-2 group">
                    <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse group-hover:scale-125 transition-transform duration-300"></div>
                    <span class="group-hover:text-green-200 transition-colors duration-300">5+ years of development</span>
                </div>
                <div class="flex items-center gap-2 group">
                    <div class="w-2 h-2 bg-blue-400 rounded-full animate-pulse delay-300 group-hover:scale-125 transition-transform duration-300"></div>
                    <span class="group-hover:text-blue-200 transition-colors duration-300">Community-driven</span>
                </div>
                <div class="flex items-center gap-2 group">
                    <div class="w-2 h-2 bg-purple-400 rounded-full animate-pulse delay-700 group-hover:scale-125 transition-transform duration-300"></div>
                    <span class="group-hover:text-purple-200 transition-colors duration-300">100% Open Source</span>
                </div>
            </div>
            
            <!-- Feature highlights -->
            <div class="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8 max-w-4xl mx-auto">
        <div class="group bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10 hover:border-white/20 transition-all duration-300 hover:scale-105 hover:bg-white/10">
            <div class="w-12 h-12 mx-auto bg-blue-400/20 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-6 h-6 text-blue-300 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                </svg>
            </div>
            <h3 class="text-lg font-semibold mb-2 group-hover:text-blue-200 transition-colors duration-300">Lightning Fast</h3>
            <p class="text-sm opacity-80 group-hover:opacity-100 transition-opacity duration-300">Built for speed and efficiency</p>
        </div>
        
        <div class="group bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10 hover:border-white/20 transition-all duration-300 hover:scale-105 hover:bg-white/10">
            <div class="w-12 h-12 mx-auto bg-purple-400/20 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-6 h-6 text-purple-300 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" xmlns="http://www.w3.org/2000/svg">
                    <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
                    <path d="M16 3.128a4 4 0 0 1 0 7.744"/>
                    <path d="M22 21v-2a4 4 0 0 0-3-3.87"/>
                    <circle cx="9" cy="7" r="4"/>
                </svg>
            </div>
            <h3 class="text-lg font-semibold mb-2 group-hover:text-purple-200 transition-colors duration-300">Team Focused</h3>
            <p class="text-sm opacity-80 group-hover:opacity-100 transition-opacity duration-300">Designed for collaborative teams</p>
        </div>
        
        <div class="group bg-white/5 backdrop-blur-sm rounded-xl p-6 border border-white/10 hover:border-white/20 transition-all duration-300 hover:scale-105 hover:bg-white/10">
            <div class="w-12 h-12 mx-auto bg-cyan-400/20 rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-6 h-6 text-cyan-300 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
                </svg>
            </div>
            <h3 class="text-lg font-semibold mb-2 group-hover:text-cyan-200 transition-colors duration-300">Secure & Reliable</h3>
            <p class="text-sm opacity-80 group-hover:opacity-100 transition-opacity duration-300">Enterprise-grade security</p>
        </div>
    </div>
        </div>
    </div>
</div>

<PageLayout>
    {@html htmlContent}
</PageLayout>
