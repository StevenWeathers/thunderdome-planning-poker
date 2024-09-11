<script>
  import { appRoutes } from '../../config';
</script>

<section class="relative bg-yellow-thunder text-gray-800 py-12">
  <div class="relative mx-auto max-w-2xl px-6 text-center lg:max-w-7xl lg:px-8">
    <h1 class="text-4xl font-bold tracking-tight sm:text-6xl">
      <span class="block lg:inline">Story Mapping with Thunderdome</span>
    </h1>
    <p class="mt-4 text-xl text-gray-700">
      Visualize your product backlog and create better user experiences with our
      powerful story mapping tool
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
        Unleash the Power of Story Mapping
      </h2>
    </div>
    <p class="mb-4">
      Story mapping is a collaborative practice that helps teams visualize the
      user journey, identify gaps in the product backlog, and prioritize
      development efforts. With Thunderdome's Story Mapping tool, teams can:
    </p>
    <ul class="list-disc pl-8 mb-6 space-y-2">
      <li>Create a shared understanding of the user experience</li>
      <li>Identify and prioritize features based on user value</li>
      <li>Plan releases and sprints more effectively</li>
      <li>Facilitate better communication between stakeholders</li>
      <li>Continuously refine and evolve the product vision</li>
    </ul>
    <p>
      Thunderdome's Story Mapping tool transforms this essential practice into
      an intuitive, collaborative, and actionable experience for product teams
      of all sizes and industries.
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
            Intuitive Visual Interface
          </h3>
        </div>
        <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">
          <li>Drag-and-drop interface for easy story organization</li>
          <li>
            Customizable columns and swimlanes to represent user activities and
            iterations
          </li>
          <li>Color-coding and tagging for improved story categorization</li>
        </ul>
      </div>
      <div class="md:w-1/2">
        <img
          src="{appConfig.PathPrefix}/img/previews/storyboard.png"
          alt="Intuitive Visual Interface Screenshot"
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
            Collaborative Mapping
          </h3>
        </div>
        <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">
          <li>Real-time collaboration for distributed teams</li>
          <li>Multi-user editing with instant updates</li>
          <li>In-context commenting and discussions</li>
        </ul>
      </div>
      <div class="md:w-1/2">
        <img
          src="{appConfig.PathPrefix}/img/previews/storyboard_users.png"
          alt="Collaborative Mapping Screenshot"
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
            Flexible Story Management
          </h3>
        </div>
        <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">
          <li>Create, edit, and prioritize user stories within the map</li>
          <li>Add descriptions, acceptance criteria, and story points</li>
        </ul>
      </div>
      <div class="md:w-1/2">
        <img
          src="{appConfig.PathPrefix}/img/previews/storyboard_story.png"
          alt="Flexible Story Management Screenshot"
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
    <!--                        Release Planning-->
    <!--                    </h3>-->
    <!--                </div>-->
    <!--                <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">-->
    <!--                    <li>Define and visualize release boundaries on the story map</li>-->
    <!--                    <li>Easily reorder and reprioritize stories across releases</li>-->
    <!--                    <li>Track progress and completion status for each release</li>-->
    <!--                    <li>Generate release notes and summaries automatically</li>-->
    <!--                </ul>-->
    <!--            </div>-->
    <!--            <div class="md:w-1/2">-->
    <!--                <img-->
    <!--                        src="/api/placeholder/800/600"-->
    <!--                        alt="Release Planning Screenshot"-->
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
    <!--                        Integration & Export-->
    <!--                    </h3>-->
    <!--                </div>-->
    <!--                <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">-->
    <!--                    <li>-->
    <!--                        Seamless integration with Thunderdome's other project management-->
    <!--                        tools-->
    <!--                    </li>-->
    <!--                    <li>-->
    <!--                        Import and export stories from popular agile management platforms-->
    <!--                    </li>-->
    <!--                    <li>-->
    <!--                        Export story maps as images or interactive HTML for easy sharing-->
    <!--                    </li>-->
    <!--                    <li>API access for custom integrations with your workflow</li>-->
    <!--                </ul>-->
    <!--            </div>-->
    <!--            <div class="md:w-1/2">-->
    <!--                <img-->
    <!--                        src="/api/placeholder/800/600"-->
    <!--                        alt="Integration & Export Screenshot"-->
    <!--                        class="rounded-lg shadow-lg"-->
    <!--                />-->
    <!--            </div>-->
    <!--        </section>-->
  </div>
</div>

<div class="bg-indigo-600 py-12">
  <div class="mx-auto max-w-7xl px-6 lg:px-8 text-center">
    <h2 class="text-3xl font-bold tracking-tight text-white sm:text-4xl">
      Ready to elevate your product planning?
    </h2>
    <p class="mt-6 text-lg leading-8 text-indigo-100">
      Join innovative product teams using Thunderdome's Story Mapping tool to
      create better user experiences and deliver value faster.
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
