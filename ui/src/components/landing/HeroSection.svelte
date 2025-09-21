<script>
  import { appRoutes } from '../../config';

  let {
    title,
    subtitle,
    primaryCta = 'Get Started',
    secondaryCta = 'Learn More',
    primaryHref = appRoutes.register,
    secondaryHref = '#features',
  } = $props();
</script>

<section class="bg-gradient-to-r from-blue-600 to-indigo-700 text-white">
  <div class="container mx-auto px-4 py-12">
    <div class="max-w-4xl mx-auto text-center">
      <h1 class="text-4xl sm:text-5xl lg:text-6xl font-bold mb-6 leading-tight">
        {@html title}
      </h1>
      <p class="text-xl sm:text-2xl text-blue-100 mb-8">
        {subtitle}
      </p>
      <div class="flex flex-col sm:flex-row items-center justify-center space-y-4 sm:space-y-0 sm:space-x-6">
        <a
          href={primaryHref}
          class="bg-white text-indigo-700 hover:bg-gray-100 font-semibold py-3 px-8 rounded-full transition duration-300 shadow-lg"
        >
          {primaryCta}
        </a>
        <a
          href={secondaryHref}
          class="bg-transparent text-white hover:bg-white/10 border border-white font-semibold py-3 px-8 rounded-full transition duration-300"
        >
          {secondaryCta}
        </a>
      </div>
    </div>
  </div>
</section>
