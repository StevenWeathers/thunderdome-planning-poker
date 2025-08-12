<script>
  import { onMount } from 'svelte';
  import HeroSection from '../../components/landing/HeroSection.svelte';
  import ProblemSolutionSection from '../../components/landing/ProblemSolutionSection.svelte';
  import FeatureSection from '../../components/landing/FeatureSection.svelte';
  import BenefitsSection from '../../components/landing/BenefitsSection.svelte';
  import FinalCTA from '../../components/landing/FinalCTA.svelte';

  onMount(() => window.scrollTo(0, 0));

  const heroData = {
    title: `Turn Retrospectives into Real Progress with <span class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-thunder to-orange-500">Thunderdome</span>`,
    subtitle: "Stop going through the motions. Start building teams that learn, adapt, and continuously improve with every sprint.",
    primaryCta: "Start Your First Retro",
    secondaryCta: "See It In Action"
  };

  const problemSolutionData = {
    title: "Why Most Retrospectives Don't Work",
    subtitle: "And how to build a culture of continuous improvement instead.",
    problemTitle: "Typical Retro Problems",
    problemItems: [
      "Same issues discussed every sprint",
      "Action items that never get done",
      "Uncomfortable silence or blame",
      "Generic formats that feel stale"
    ],
    solutionTitle: "With Thunderdome Retrospectives",
    solutionItems: [
      "Fresh insights every retrospective",
      "Trackable action items with follow-through",
      "Safe space for honest feedback",
      "Dynamic formats that energize teams"
    ],
    conclusion: "<strong>The difference?</strong> Retrospectives become the engine of improvement, not just another meeting on the calendar."
  };

  const features = [
    {
      title: "Formats That Actually Engage",
      description: "Move beyond boring retrospectives. Use proven formats that spark meaningful discussions and fresh perspectives.",
      features: [
        "<strong>Battle-tested templates:</strong> Start/Stop/Continue, 4Ls, Mad/Sad/Glad, and more",
        "<strong>Create custom templates:</strong> Design and save your own retrospective formats tailored to your team's needs",
        "<strong>Mix it up easily:</strong> Switch formats to keep retrospectives energizing and insightful",
        "<strong>Flexible structure:</strong> Adapt any template to your team's unique challenges and wins"
      ],
      imageDark: "/img/previews/retro_template_preview.png",
      imageLight: "/img/previews/retro_template_light_preview.png",
      imageAlt: "Multiple retrospective formats and templates"
    },
    {
      title: "Everyone's Voice Gets Heard",
      description: "Create psychological safety where introverts contribute as much as extroverts, and honest feedback flows freely.",
      features: [
        "<strong>Remote-friendly participation:</strong> Equal engagement whether you're in-person or distributed",
        "<strong>Anonymous feedback options:</strong> Share sensitive insights without fear of judgment",
        "<strong>Democratic prioritization:</strong> Vote on what matters most to focus discussion time",
        "<strong>Guided facilitation:</strong> Built-in tools to keep discussions productive and inclusive"
      ],
      imageDark: "/img/previews/retro.png",
      imageLight: "/img/previews/retro_light.png",
      imageAlt: "Real-time collaborative retrospective session"
    },
    {
      title: "Follow Through, Not Follow Up",
      description: "Turn insights into action with accountability that actually works. No more forgotten commitments or repeated discussions.",
      features: [
        "<strong>Seamless action creation:</strong> Convert discussion points into trackable commitments instantly",
        "<strong>Cross-sprint visibility:</strong> See progress on improvements from previous retrospectives"
      ],
      imageDark: "/img/previews/action_items.png",
      imageLight: "/img/previews/action_items_light.png",
      imageAlt: "Action item tracking and progress monitoring"
    }
  ];

  const benefitsData = {
    title: "Transform Your Team's Continuous Improvement",
    description: "Agile teams using Thunderdome retrospectives report stronger team bonds, faster problem-solving, and measurable process improvements.",
    benefits: [
      {
        title: "Psychological Safety",
        description: "Create an environment where team members feel safe to share honest feedback"
      },
      {
        title: "Actionable Insights",
        description: "Move beyond venting to identifying concrete improvements you can implement"
      },
      {
        title: "Sustained Change",
        description: "Track improvements over time and build momentum for lasting team growth"
      }
    ]
  };

  const finalCtaData = {
    title: "Ready to Make Every Retrospective Count?",
    description: "Stop wasting time on retrospectives that don't drive change. Start building a culture of continuous improvement that makes your team stronger every sprint.",
    primaryCta: "Run Your First Effective Retro",
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