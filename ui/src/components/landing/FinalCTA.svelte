<script>
  import { appRoutes } from '../../config';

  let {
    title,
    description,
    primaryCta = 'Get Started',
    secondaryCta = 'Sign In to Continue',
    primaryHref = appRoutes.register,
    secondaryHref = appRoutes.login,
  } = $props();
</script>

<div class="bg-gradient-to-r from-indigo-600 to-blue-700 py-16">
  <div class="mx-auto max-w-4xl px-6 lg:px-8 text-center">
    <h2 class="text-4xl font-bold tracking-tight text-white sm:text-5xl mb-6">
      {title}
    </h2>
    <p class="text-xl leading-8 text-indigo-100 mb-8 max-w-2xl mx-auto">
      {description}
    </p>

    <div class="flex flex-col sm:flex-row items-center justify-center space-y-4 sm:space-y-0 sm:space-x-6">
      <a
        href={primaryHref}
        class="bg-white text-indigo-700 hover:bg-gray-100 font-semibold py-4 px-8 rounded-full transition duration-300 shadow-lg text-lg"
      >
        {primaryCta}
      </a>
      <a
        href={secondaryHref}
        class="bg-transparent text-white hover:bg-white/10 border-2 border-white font-semibold py-4 px-8 rounded-full transition duration-300 text-lg"
      >
        {secondaryCta}
      </a>
    </div>
  </div>
</div>
