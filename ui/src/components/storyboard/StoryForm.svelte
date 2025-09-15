<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import Editor from '../forms/Editor.svelte';
  import { User, MessageCircleMore, ChevronRight, ChevronDown } from 'lucide-svelte';
  import { onMount } from 'svelte';
  
  import type { NotificationService } from '../../types/notifications';

  interface Props {
    toggleStoryForm?: any;
    sendSocketEvent?: any;
    notifications: NotificationService;
    story?: any;
    colorLegend?: any;
    users?: any;
  }

  let {
    toggleStoryForm = () => {},
    sendSocketEvent = () => {},
    notifications,
    story = $bindable({}),
    colorLegend = [],
    users = []
  }: Props = $props();

  const isAbsolute = new RegExp('^([a-z]+://|//)', 'i');

  let userComment = $state('');
  let selectedComment = $state(null);
  let selectedCommentContent = $state('');
  let actionsHidden = $state(true);
  let discussionHidden = $state(true);
  let focusInput: any = $state();

  let userMap = $derived(users.reduce((prev, usr) => {
    prev[usr.id] = usr.name;
    return prev;
  }, {}));

  function handleStoryDelete() {
    sendSocketEvent('delete_story', story.id);
    toggleStoryForm();
  }

  function markClosed() {
    sendSocketEvent(
      'update_story_closed',
      JSON.stringify({
        storyId: story.id,
        closed: true,
      }),
    );
  }

  function markOpen() {
    sendSocketEvent(
      'update_story_closed',
      JSON.stringify({
        storyId: story.id,
        closed: false,
      }),
    );
  }

  const changeColor = color => () => {
    sendSocketEvent(
      'update_story_color',
      JSON.stringify({
        storyId: story.id,
        color,
      }),
    );
  };

  const updateName = evt => {
    const name = evt.target.value;
    sendSocketEvent(
      'update_story_name',
      JSON.stringify({
        storyId: story.id,
        name,
      }),
    );
  };

  const updateContent = () => {
    sendSocketEvent(
      'update_story_content',
      JSON.stringify({
        storyId: story.id,
        content: story.content,
      }),
    );
  };

  const updatePoints = evt => {
    const points = parseInt(evt.target.value, 10);
    sendSocketEvent(
      'update_story_points',
      JSON.stringify({
        storyId: story.id,
        points,
      }),
    );
  };

  const updateLink = evt => {
    const link = evt.target.value;
    if (link !== '' && !isAbsolute.test(link)) {
      notifications.danger('Link must be an absolute URL');
      return;
    }

    sendSocketEvent(
      'update_story_link',
      JSON.stringify({
        storyId: story.id,
        link,
      }),
    );
  };

  const handleCommentSubmit = () => {
    if (userComment !== '') {
      sendSocketEvent(
        'add_story_comment',
        JSON.stringify({ storyId: story.id, comment: userComment }),
      );
      userComment = '';
    }
  };

  const toggleCommentEdit = comment => () => {
    selectedComment = comment;
    if (comment !== null) {
      selectedCommentContent = comment.comment;
    }
  };

  const handleCommentEdit = () => {
    sendSocketEvent(
      'edit_story_comment',
      JSON.stringify({
        commentId: selectedComment.id,
        comment: selectedCommentContent,
      }),
    );
    selectedComment = null;
    selectedCommentContent = '';
  };

  const handleCommentDelete = (commentId: String) => () => {
    sendSocketEvent('delete_story_comment', JSON.stringify({ commentId }));
  };

  const toggleMoreActions = () => {
    actionsHidden = !actionsHidden;
  }

  const toggleDiscussion = () => {
    discussionHidden = !discussionHidden;
  }

  onMount(() => {
    focusInput?.focus();
  });
</script>

<style>
  .colorcard-gray {
    @apply bg-gray-400;
    @apply ring-gray-400;
  }

  .colorcard-gray:hover {
    @apply bg-gray-600;
    @apply ring-gray-600;
  }

  .colorcard-red {
    @apply bg-red-400;
    @apply ring-red-400;
  }

  .colorcard-red:hover {
    @apply bg-red-600;
    @apply ring-red-600;
  }

  .colorcard-orange {
    @apply bg-orange-400;
    @apply ring-orange-400;
  }

  .colorcard-orange:hover {
    @apply bg-orange-600;
    @apply ring-orange-600;
  }

  .colorcard-yellow {
    @apply bg-yellow-400;
    @apply ring-yellow-400;
  }

  .colorcard-yellow:hover {
    @apply bg-yellow-600;
    @apply ring-yellow-600;
  }

  .colorcard-green {
    @apply bg-green-400;
    @apply ring-green-400;
  }

  .colorcard-green:hover {
    @apply bg-green-600;
    @apply ring-green-600;
  }

  .colorcard-teal {
    @apply bg-teal-400;
    @apply ring-teal-400;
  }

  .colorcard-teal:hover {
    @apply bg-teal-600;
    @apply ring-teal-600;
  }

  .colorcard-blue {
    @apply bg-blue-400;
    @apply ring-blue-400;
  }

  .colorcard-blue:hover {
    @apply bg-blue-600;
    @apply ring-blue-600;
  }

  .colorcard-indigo {
    @apply bg-indigo-400;
    @apply ring-indigo-400;
  }

  .colorcard-indigo:hover {
    @apply bg-indigo-600;
    @apply ring-indigo-600;
  }

  .colorcard-purple {
    @apply bg-purple-400;
    @apply ring-purple-400;
  }

  .colorcard-purple:hover {
    @apply bg-purple-600;
    @apply ring-purple-600;
  }

  .colorcard-pink {
    @apply bg-pink-400;
    @apply ring-pink-400;
  }

  .colorcard-pink:hover {
    @apply bg-pink-600;
    @apply ring-pink-600;
  }
</style>

<Modal closeModal={toggleStoryForm} widthClasses="w-full md:w-3/4 xl:1/2 2xl:w-2/5" ariaLabel={$LL.modalStoryboardStory()}>
  <div class="p-6 overflow-y-auto h-full">
      <div class="space-y-6">
          <!-- Story Name (Full Width) -->
          <div>
              <label for="storyName" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
                  Story Name <span class="text-red-500">*</span>
              </label>
              <!-- <input 
                  type="text" 
                  id="story-name"
                  placeholder="e.g. User Authentication System"
                  class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                  value=""
              > -->
              <TextInput
                id="storyName"
                onchange={updateName}
                value={story.name}
                placeholder="e.g. User Authentication System"
                name="storyName"
                bind:this={focusInput}
              />
          </div>

          <!-- Story Points and Color -->
          <div>
              <div class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">Story Color</div>
              <div class="flex space-x-2 pt-1">
              {#each colorLegend as color}
                <button
                  onclick={changeColor(color.color)}
                  class="w-8 h-8 rounded-full colorcard-{color.color}
                  hover:scale-110 transition-transform dark:ring-offset-gray-800 {story.color ===
                  color.color
                    ? `ring-2 ring-offset-2`
                    : ''}"
                  title="{color.color}{color.legend !== ''
                    ? ` - ${color.legend}`
                    : ''}"><span class="hidden">change color</span></button>
              {/each}
              </div>
          </div>

          <!-- Story Content -->
          <div>
              <label for="story-content" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
                  Story Description
              </label>
              <div class="border border-gray-300 dark:border-gray-600 rounded-lg overflow-hidden">
                  <!-- <textarea 
                      id="story-content"
                      rows="8"
                      placeholder="As a [user type], I want [functionality] so that [benefit]..."
                      class="w-full px-4 py-3 border-0 dark:bg-gray-700 dark:text-white focus:ring-0 resize-none"
                  ></textarea> -->
                  <div class="bg-white">
                  <Editor
                    content={story.content}
                    placeholder="Enter story content"
                    id="storyDescription"
                    handleTextChange={c => {
                      story.content = c;
                      updateContent();
                    }}
                  />
                  </div>
              </div>
          </div>

          <!-- More Actions Section (Collapsed by default) -->
          <div class="border-t border-gray-200 dark:border-gray-700 pt-4">
              <button 
                  id="more-actions-toggle"
                  class="flex items-center space-x-2 font-bold text-lg lg:text-xl text-gray-800 dark:text-gray-200 hover:text-blue-800 dark:hover:text-blue-200 transition-colors"
                  onclick={toggleMoreActions}

              >
                  {#if actionsHidden}
                    <ChevronRight class="inline-block w-5 h-5" />
                  {:else}
                    <ChevronDown class="inline-block w-5 h-5" />
                  {/if}
                  <span>Additional Details</span>
              </button>
              
              <div id="more-actions-content" class="{actionsHidden ? 'hidden' : ''} mt-4 space-y-4 ps-6 border-l-2 border-gray-100 dark:border-gray-700">
                  <!-- Story Points -->
                  <div>
                      <label for="storyPoints" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
                          Story Points
                      </label>
                      <!-- <select class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors">
                          <option value="0">0</option>
                          <option value="1">1</option>
                          <option value="2">2</option>
                          <option value="3">3</option>
                          <option value="5">5</option>
                          <option value="8">8</option>
                          <option value="13">13</option>
                          <option value="21">21</option>
                      </select> -->
                      <input
                          class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                          id="storyPoints"
                          type="number"
                          min="0"
                          max="999"
                          bind:value="{story.points}"
                          onchange={updatePoints}
                          placeholder="Enter story points e.g. 1, 2, 3, 5,
                                        8"
                          name="storyPoints"
                        />
                  </div>
                  
                  <!-- Story Link -->
                  <div>
                      <label for="storyLink" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
                          Story Link
                          <span class="text-gray-500 dark:text-gray-400 font-normal ms-1 text-lg">(Optional)</span>
                      </label>
                      <div class="relative">
                        <TextInput
                          id="storyLink"
                          onchange={updateLink}
                          value={story.link}
                          placeholder="https://jira.company.com/browse/PROJ-123"
                          name="storyLink"
                          type="url"
                          class="ps-10"
                        />
                          <!-- <input 
                              class="w-full px-4 py-3 ps-10 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                          > -->
                          <svg class="w-5 h-5 text-gray-400 dark:text-gray-500 absolute left-3 top-1/2 transform -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
                          </svg>
                      </div>
                  </div>

                  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <!-- Story Status-->
                    <div>
                      <div class="text-gray-700 dark:text-gray-300 mb-3 text-lg">
                        Story Status
                      </div>
                      {#if !story.closed}
                        <HollowButton color="orange" onClick={markClosed}>
                          <svg class="w-4 h-4 inline me-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                          </svg> Mark as Closed
                        </HollowButton>
                      {:else}
                        <HollowButton color="green" onClick={markOpen}
                          >Reopen story
                        </HollowButton>
                      {/if}
                    </div>

                    <!-- Delete Story-->
                    <div>
                      <h4 class="text-red-800 dark:text-red-400 mb-3 flex items-center text-lg">
                            <svg class="w-5 h-5 me-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.966-.833-2.736 0L3.478 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                            </svg>
                            Danger Zone
                      </h4>
                      <div class="flex space-x-3">
                          <HollowButton color="red" onClick={handleStoryDelete}>
                            <svg class="w-4 h-4 inline me-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1-1H8a1 1 0 00-1 1v3M4 7h16"></path>
                              </svg> Delete Story
                          </HollowButton>
                      </div>
                    </div>
                  </div>
              </div>
          </div>

          <!-- Discussion Section -->
          <div class="border-t border-gray-200 dark:border-gray-700 pt-4 lg:pt-6">
              <h3 class="text-lg lg:text-xl text-gray-900 dark:text-white flex items-center mb-4">
                  <button 
                    id="discussion-toggle"
                    class="flex items-center space-x-2 font-bold text-gray-800 dark:text-gray-200 hover:text-blue-800 dark:hover:text-blue-200 transition-colors"
                    onclick={toggleDiscussion}

                >
                    {#if discussionHidden}
                      <ChevronRight class="inline-block w-5 h-5" />
                    {:else}
                      <ChevronDown class="inline-block w-5 h-5" />
                    {/if}
                    <span>Discussion</span><MessageCircleMore class="inline-block ms-2 w-5 h-5" />
                </button>
                
                  
                  <span class="text-normal ms-auto bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 leading-none py-2 {story.comments.length > 9 ? 'px-2' : 'px-3'} rounded-full">{story.comments ? `${story.comments.length}` : ''}</span>
              </h3>

              <div class="{discussionHidden ? 'hidden' : ''}">
              <!-- Comments List -->
               {#if story.comments}
                <div class="space-y-4 mb-6">
                    {#each story.comments as comment}
                      <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4 shadow-sm" data-commentid="{comment.id}">
                          <div class="flex items-center space-x-2 mb-2">
                              <div class="w-6 h-6 bg-blue-500 rounded-full flex items-center justify-center text-white text-xs font-medium">
                                  <User class="h-4 w-4 inline-block" />
                              </div>
                              <span class="text-gray-900 dark:text-white">{userMap[
                            comment.user_id
                          ]}</span>
                              <!-- <span class="text-xs text-gray-500 dark:text-gray-400">2 min ago</span> -->
                          </div>
                          {#if selectedComment !== null && selectedComment.id === comment.id}
                            <div class="w-full my-2">
                              <textarea
                                class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                                      rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                                      focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
                                bind:value="{selectedCommentContent}"></textarea>
                              <div class="text-right">
                                <HollowButton
                                  color="blue"
                                  onClick={toggleCommentEdit(null)}
                                >
                                  {$LL.cancel()}
                                </HollowButton>
                                <HollowButton
                                  color="green"
                                  onClick={handleCommentEdit}
                                  disabled={selectedCommentContent === ''}
                                >
                                  {$LL.updateComment()}
                                </HollowButton>
                              </div>
                            </div>
                          {:else}
                            <p class="text-gray-700 dark:text-gray-300">{comment.comment}</p>
                          {/if}

                          {#if comment.user_id === $user.id && !(selectedComment !== null && selectedComment.id === comment.id)}
                          <div class="flex justify-end space-x-2 mt-2">
                              <button class="text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300" onclick={toggleCommentEdit(comment)}>{$LL.edit()}</button>
                              <button class="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300" onclick={handleCommentDelete(comment.id)}>{$LL.delete()}</button>
                          </div>
                          {/if}                        
                      </div>
                    {/each}
                </div>
              {/if}

              <!-- Add Comment -->
              <div>
                  <textarea 
                      placeholder="Write a comment..."
                      rows="3"
                      class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none mb-3"
                      bind:value="{userComment}"
                  ></textarea>
                  <div class="flex justify-end">
                    <button class="bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors font-medium"
                      onclick={handleCommentSubmit}
                      disabled={userComment === ''}>
                        Post Comment
                    </button>
                  </div>
              </div>
            </div>
          </div>
      </div>
  </div>
</Modal>
