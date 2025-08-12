<script>
  import { onMount } from 'svelte';
  import HeroSection from '../../components/landing/HeroSection.svelte';
  import ProblemSolutionSection from '../../components/landing/ProblemSolutionSection.svelte';
  import FeatureSection from '../../components/landing/FeatureSection.svelte';
  import BenefitsSection from '../../components/landing/BenefitsSection.svelte';
  import FinalCTA from '../../components/landing/FinalCTA.svelte';

  onMount(() => window.scrollTo(0, 0));

  const heroData = {
    title: `Turn Estimation Chaos into Team Consensus with <span class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-thunder to-orange-500">Thunderdome</span> Planning Poker`,
    subtitle: "Transform wild guesses and estimation bias into data-driven story points. Build accurate backlogs that help your team plan realistic sprints.",
    primaryCta: "Start Estimating Smart",
    secondaryCta: "See How It Works"
  };

  const problemSolutionData = {
    title: "Why Traditional Estimation Fails Teams",
    subtitle: "Move from wild guesses and endless debates to data-driven consensus that actually works.",
    problemTitle: "Traditional Estimation Problems",
    problemItems: [
      "Loudest voice wins the estimate",
      "Anchoring bias skews all estimates", 
      "Unrealistic commitments lead to burnout",
      "Team members afraid to disagree"
    ],
    solutionTitle: "With Thunderdome Planning Poker",
    solutionItems: [
      "Every team member's input counts equally",
      "Independent thinking before consensus",
      "Realistic estimates based on team wisdom",
      "Safe space for honest complexity assessment"
    ],
    conclusion: "<strong>The result?</strong> Accurate story points that help your team make realistic sprint commitments and improve velocity predictions over time."
  };

  const features = [
    {
      title: "Scales That Match Your Reality",
      description: "Stop forcing your team into estimation scales that don't fit. Use proven scales or create your own that reflect how your team actually thinks about work.",
      features: [
        "<strong>Classic scales ready to go:</strong> Fibonacci, T-shirt sizes, and other proven estimation approaches",
        "<strong>Build custom scales:</strong> Create estimation ranges that match your team's workflow and complexity levels",
        "<strong>Special situation cards:</strong> Handle unknowns with \"?\" cards or infinite complexity with \"∞\" options",
        "<strong>Switch scales easily:</strong> Use different approaches for different project types or team preferences"
      ],
      imageDark: "/img/previews/estimation_scales.png",
      imageLight: "/img/previews/estimation_scales_light.png",
      imageAlt: "Customizable estimation scales including Fibonacci and custom options"
    },
    {
      title: "Estimate Together, Anywhere",
      description: "Whether your team is in the same room or spread across time zones, everyone participates equally in the estimation process.",
      features: [
        "<strong>Real-time synchronization:</strong> See who's voted and when, with instant updates across all devices",
        "<strong>Seamless remote participation:</strong> No difference between in-person and distributed team members",
        "<strong>Multi-user sessions:</strong> Handle large teams without technical hiccups or coordination issues"
      ],
      imageDark: "/img/previews/voting_interactive.png",
      imageLight: "/img/previews/voting_interactive_light.png",
      imageAlt: "Real-time collaborative planning poker session"
    },
    {
      title: "Unbiased Estimation Made Easy",
      description: "Eliminate estimation bias and anchoring effects. Get the true wisdom of your team, not just the opinion of whoever speaks first.",
      features: [
        "<strong>Anonymous voting options:</strong> Remove personality and hierarchy bias from the estimation process",
        "<strong>Simultaneous reveals:</strong> Everyone thinks independently before seeing others' estimates",
        "<strong>Outlier identification:</strong> Quickly spot estimation differences that need discussion and alignment"
      ],
      imageDark: "/img/previews/voting_anon.png",
      imageLight: "/img/previews/voting_anon_light.png",
      imageAlt: "Anonymous voting interface showing unbiased estimation process"
    }
  ];

  const benefitsData = {
    title: "Build Accurate Story Points Your Team Can Trust",
    description: "Teams using Thunderdome Planning Poker report more consistent story pointing, better velocity tracking, and improved sprint planning confidence.",
    benefits: [
      {
        title: "Story Point Accuracy",
        description: "Leverage collective team intelligence for more reliable and consistent story point estimates"
      },
      {
        title: "Team Consensus", 
        description: "Build shared understanding of story complexity and effort before pointing stories"
      },
      {
        title: "Velocity Tracking",
        description: "Improve sprint planning with consistent story points that reflect true team capacity"
      }
    ]
  };

  const finalCtaData = {
    title: "Ready to Point Stories Like a Pro?",
    description: "Stop guessing story complexity and start measuring it accurately. Build consistent story points that improve your team's velocity tracking and sprint planning.",
    primaryCta: "Start Pointing Stories",
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