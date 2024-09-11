<script>
  import { appRoutes } from '../../config';
</script>

<section class="relative bg-yellow-thunder text-gray-800 py-12">
  <div class="relative mx-auto max-w-2xl px-6 text-center lg:max-w-7xl lg:px-8">
    <h1 class="text-4xl font-bold tracking-tight sm:text-6xl">
      <span class="block lg:inline">Agile Poker Planning with Thunderdome</span>
    </h1>
    <p class="mt-4 text-xl text-gray-700">
      Streamline your sprint planning process with our intuitive and powerful
      poker planning tool.
    </p>
  </div>
</section>

<div class="py-16 mx-auto max-w-2xl px-6 lg:max-w-7xl lg:px-8">
  <section class="mb-20 text-gray-700 dark:text-gray-300">
    <div class="mb-4">
      <div class="title-line bg-yellow-thunder"></div>
      <h2
        class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
      >
        Elevate Your Estimations
      </h2>
    </div>
    <p class="mb-4">
      Agile Poker Planning is a consensus-based estimation technique used in
      agile software development. It's a powerful method that helps teams
      collaboratively estimate the effort required for user stories or tasks. By
      using Thunderdome's Poker Planning tool, teams can:
    </p>
    <ul class="list-disc pl-8 mb-6 space-y-2">
      <li>Encourage team participation and engagement</li>
      <li>Reduce influence bias in estimations</li>
      <li>Promote discussion and alignment on project scope</li>
      <li>Improve estimation accuracy over time</li>
      <li>Streamline the sprint planning process</li>
    </ul>
    <p>
      Thunderdome's Poker Planning tool transforms this essential practice into
      a seamless, engaging, and actionable experience for agile teams of all
      sizes and industries.
    </p>
  </section>

  <div class="space-y-24">
    <section class="flex flex-col md:flex-row items-center">
      <div class="md:w-1/2 md:pr-8 mb-8 md:mb-0">
        <div class="mb-4">
          <div class="title-line bg-yellow-thunder"></div>
          <h3
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
          >
            Custom Estimation Scales
          </h3>
        </div>
        <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">
          <li>
            Choose from predefined scales like Fibonacci or create your own
            custom estimation scales
          </li>
          <li>
            Tailor scales to fit your team's unique needs and project
            requirements
          </li>
          <li>
            Easily switch between different scales for various project types
          </li>
          <li>
            Include special cards like "?" or "∞" for handling uncertainties
          </li>
        </ul>
      </div>
      <div class="md:w-1/2">
        <img
          src="{appConfig.PathPrefix}/img/previews/estimation_scales.png"
          alt="Custom Estimation Scales Screenshot"
          class="rounded-lg shadow-lg max-h-64 max-w-64 mx-auto"
        />
      </div>
    </section>

    <section class="flex flex-col md:flex-row-reverse items-center">
      <div class="md:w-1/2 md:pl-8 mb-8 md:mb-0">
        <div class="mb-4">
          <div class="title-line bg-yellow-thunder"></div>
          <h3
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
          >
            Real-time Collaboration
          </h3>
        </div>
        <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">
          <li>
            Conduct poker planning sessions with your team in real-time,
            regardless of location
          </li>
          <li>Seamless multi-user participation with instant updates</li>
          <li>Suitable for both co-located and distributed agile teams</li>
        </ul>
      </div>
      <div class="md:w-1/2">
        <img
          src="{appConfig.PathPrefix}/img/previews/voting_interactive.png"
          alt="Real-time Collaboration Screenshot"
          class="rounded-lg shadow-lg max-h-64 max-w-64 mx-auto"
        />
      </div>
    </section>

    <section class="flex flex-col md:flex-row items-center">
      <div class="md:w-1/2 md:pr-8 mb-8 md:mb-0">
        <div class="mb-4">
          <div class="title-line bg-yellow-thunder"></div>
          <h3
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
          >
            Voting Visibility Control
          </h3>
        </div>
        <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">
          <li>
            Choose whether to show or hide voter identities to reduce bias
          </li>
          <li>
            Reveal all estimates simultaneously to encourage independent
            thinking
          </li>
          <!--                    <li>-->
          <!--                        Option to keep voting open or set time limits for each estimation-->
          <!--                        round-->
          <!--                    </li>-->
          <li>
            Easily identify outliers and facilitate discussions on estimate
            differences
          </li>
        </ul>
      </div>
      <div class="md:w-1/2">
        <img
          src="{appConfig.PathPrefix}/img/previews/voting_anon.png"
          alt="Voting Visibility Control Screenshot"
          class="rounded-lg shadow-lg max-h-64 max-w-64 mx-auto"
        />
      </div>
    </section>

    <!--        <section class="flex flex-col md:flex-row-reverse items-center">-->
    <!--            <div class="md:w-1/2 md:pl-8 mb-8 md:mb-0">-->
    <!--                <div class="mb-4">-->
    <!--                    <div class="title-line bg-yellow-thunder"></div>-->
    <!--                    <h3-->
    <!--                            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"-->
    <!--                    >-->
    <!--                        Integration with Projects-->
    <!--                    </h3>-->
    <!--                </div>-->
    <!--                <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">-->
    <!--                    <li>-->
    <!--                        Seamlessly integrate poker planning sessions with your existing-->
    <!--                        projects in Thunderdome-->
    <!--                    </li>-->
    <!--                    <li>Import user stories directly from your project backlog</li>-->
    <!--                    <li>-->
    <!--                        Automatically update story points in your project management tool-->
    <!--                    </li>-->
    <!--                    <li>-->
    <!--                        Link estimation results to sprints for better planning and tracking-->
    <!--                    </li>-->
    <!--                </ul>-->
    <!--            </div>-->
    <!--            <div class="md:w-1/2">-->
    <!--                <img-->
    <!--                        src="/api/placeholder/800/600"-->
    <!--                        alt="Project Integration Screenshot"-->
    <!--                        class="rounded-lg shadow-lg"-->
    <!--                />-->
    <!--            </div>-->
    <!--        </section>-->

    <!--        <section class="flex flex-col md:flex-row items-center">-->
    <!--            <div class="md:w-1/2 md:pr-8 mb-8 md:mb-0">-->
    <!--                <div class="mb-4">-->
    <!--                    <div class="title-line bg-yellow-thunder"></div>-->
    <!--                    <h3-->
    <!--                            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"-->
    <!--                    >-->
    <!--                        Historical Data & Analytics-->
    <!--                    </h3>-->
    <!--                </div>-->
    <!--                <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">-->
    <!--                    <li>-->
    <!--                        Track and analyze estimation trends over time to improve accuracy-->
    <!--                    </li>-->
    <!--                    <li>Generate reports on team velocity and estimation consistency</li>-->
    <!--                    <li>-->
    <!--                        Compare estimated vs. actual effort to refine future estimations-->
    <!--                    </li>-->
    <!--                    <li>-->
    <!--                        Visualize data with interactive charts and graphs for easy-->
    <!--                        interpretation-->
    <!--                    </li>-->
    <!--                </ul>-->
    <!--            </div>-->
    <!--            <div class="md:w-1/2">-->
    <!--                <img-->
    <!--                        src="/api/placeholder/800/600"-->
    <!--                        alt="Historical Data & Analytics Screenshot"-->
    <!--                        class="rounded-lg shadow-lg"-->
    <!--                />-->
    <!--            </div>-->
    <!--        </section>-->
  </div>
</div>

<div class="bg-indigo-600 py-12">
  <div class="mx-auto max-w-7xl px-6 lg:px-8 text-center">
    <h2 class="text-3xl font-bold tracking-tight text-white sm:text-4xl">
      Ready to revolutionize your sprint planning?
    </h2>
    <p class="mt-6 text-lg leading-8 text-indigo-100">
      Join thousands of agile teams using Thunderdome's Poker Planning tool to
      streamline their estimation process and boost productivity.
    </p>
    <div class="mt-10 flex items-center justify-center gap-x-6">
      <div class="text-center">
        <div
          class="flex flex-col sm:flex-row items-center justify-center space-y-4 sm:space-y-0 sm:space-x-4"
        >
          <a
            href="{appRoutes.login}"
            class="bg-white bg-indigo-600 hover:bg-gray-100 font-semibold py-2 px-6 rounded-full transition duration-300"
          >
            Login
          </a>
          <div class="text-white font-medium">or</div>
          <a
            href="{appRoutes.register}"
            class="bg-transparent text-white hover:bg-white/10 border border-white font-semibold py-2 px-6 rounded-full transition duration-300"
          >
            Register
          </a>
        </div>
        <p class="text-white mt-2">to get started</p>
      </div>
    </div>
  </div>
</div>
