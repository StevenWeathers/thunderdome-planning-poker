<script>
  import { onMount } from 'svelte';
  import HeroSection from '../../components/landing/HeroSection.svelte';
  import ProblemSolutionSection from '../../components/landing/ProblemSolutionSection.svelte';
  import FeatureSection from '../../components/landing/FeatureSection.svelte';
  import BenefitsSection from '../../components/landing/BenefitsSection.svelte';
  import FinalCTA from '../../components/landing/FinalCTA.svelte';

  onMount(() => window.scrollTo(0, 0));

  const heroData = {
    title: `Transform Ideas into Impact with <span class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-thunder to-orange-500">Thunderdome</span> Story Mapping`,
    subtitle: "Turn complex product visions into clear, actionable user journeys that your entire team can rally behind.",
    primaryCta: "Start Mapping for Free",
    secondaryCta: "See How It Works"
  };

  const problemSolutionData = {
    title: "Why Story Mapping Changes Everything",
    subtitle: "Stop building features in isolation. Start building experiences that matter.",
    problemTitle: "Before Story Mapping",
    problemItems: [
      "Disconnected feature lists",
      "Unclear user journeys", 
      "Misaligned team priorities",
      "Feature confusion"
    ],
    solutionTitle: "With Thunderdome Story Mapping",
    solutionItems: [
      "Complete user journey visualization",
      "Strategic feature prioritization",
      "Unified team understanding", 
      "Clear product roadmap"
    ],
    conclusion: "<strong>The result?</strong> Products that users actually want, delivered by teams that know exactly what they're building and why."
  };

  const features = [
    {
      title: "Build Stories Visually, Not Just Verbally",
      description: "Skip the endless spreadsheets. Create story maps that actually show how your product flows from the user's perspective.",
      features: [
        "<strong>Intuitive drag-and-drop:</strong> Reorganize stories as easily as moving sticky notes",
        "<strong>Smart organization:</strong> Custom swimlanes that adapt to your team's workflow",
        "<strong>Visual clarity:</strong> Color-coding and tags that make complex backlogs instantly understandable"
      ],
      imageDark: "/img/previews/storyboard_20250812_dark.png",
      imageLight: "/img/previews/storyboard_20250812_light.png",
      imageAlt: "Visual story mapping interface showing drag-and-drop functionality"
    },
    {
      title: "Map Together, Ship Faster",
      description: "Whether your team is across the hall or across the globe, everyone stays in sync with real-time collaboration.",
      features: [
        "<strong>Live collaboration:</strong> See changes as they happen, no refresh needed",
        "<strong>Seamless multi-user editing:</strong> Multiple people, zero conflicts",
        "<strong>Contextual discussions:</strong> Have conversations right where the work happens"
      ],
      imageDark: "/img/previews/storyboard_users_20250812_dark.png",
      imageLight: "/img/previews/storyboard_users_20250812_light.png",
      imageAlt: "Team collaboration in real-time on story maps"
    },
    {
      title: "Stories That Actually Tell a Story",
      description: "Transform vague requirements into crystal-clear user stories with all the context your team needs to build with confidence.",
      features: [
        "<strong>Rich story details:</strong> Descriptions, acceptance criteria, and estimates all in one place",
        "<strong>Priority-driven organization:</strong> Focus on what matters most to your users"
      ],
      imageDark: "/img/previews/storyboard_story_20250812_dark.png",
      imageLight: "/img/previews/storyboard_story_20250812_light.png",
      imageAlt: "Detailed story management with descriptions and acceptance criteria"
    }
  ];

  const benefitsData = {
    title: "Join Teams Who've Transformed Their Product Development",
    description: "From startups to enterprise teams, product managers are using Thunderdome to turn chaotic backlogs into clear roadmaps that drive real results.",
    benefits: [
      {
        title: "User Journey Clarity",
        description: "Visualize the complete user experience from start to finish"
      },
      {
        title: "Team Alignment", 
        description: "Get everyone on the same page about priorities and user value"
      },
      {
        title: "Better Prioritization",
        description: "Focus development efforts on features that truly matter to users"
      }
    ]
  };

  const finalCtaData = {
    title: "Ready to Map Your Way to Better Products?",
    description: "Stop guessing what to build next. Start with story mapping and build products your users will love—backed by data, driven by user needs.",
    primaryCta: "Create Your First Story Map",
    secondaryCta: "Sign In to Continue"
  };
</script>

<HeroSection {...heroData} />

<div class="py-16 mx-auto max-w-2xl px-6 lg:max-w-7xl lg:px-8">
  <ProblemSolutionSection {...problemSolutionData} />

  <div id="features" class="space-y-24">
    {#each features as feature, index}
      <FeatureSection {...feature} {index} />
    {/each}
  </div>

  <BenefitsSection {...benefitsData} />
</div>

<FinalCTA {...finalCtaData} />