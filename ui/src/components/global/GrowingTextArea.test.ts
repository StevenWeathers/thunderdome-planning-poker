import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import GrowingTextArea from './GrowingTextArea.svelte';

describe('GrowingTextArea component', () => {
  it('should render successfully', () => {
    render(GrowingTextArea, {});

    const textarea = page.getByRole('textbox');
    expect(textarea).toBeTruthy();
  });

  it('should render with placeholder', () => {
    const { container } = render(GrowingTextArea, {
      placeholder: 'Enter your text here',
    });

    const textarea = container.querySelector('textarea');
    expect(textarea?.getAttribute('placeholder')).toBe('Enter your text here');
  });

  it('should render with id attribute', () => {
    const { container } = render(GrowingTextArea, {
      id: 'test-textarea',
    });

    const textarea = container.querySelector('textarea');
    expect(textarea?.getAttribute('id')).toBe('test-textarea');
  });

  it('should render with name attribute', () => {
    const { container } = render(GrowingTextArea, {
      name: 'description',
    });

    const textarea = container.querySelector('textarea');
    expect(textarea?.getAttribute('name')).toBe('description');
  });

  it('should render with required attribute when required is true', () => {
    const { container } = render(GrowingTextArea, {
      required: true,
    });

    const textarea = container.querySelector('textarea');
    expect(textarea?.getAttribute('required')).toBe('');
  });

  it('should not have required attribute when required is false', () => {
    const { container } = render(GrowingTextArea, {
      required: false,
    });

    const textarea = container.querySelector('textarea');
    expect(textarea?.hasAttribute('required')).toBe(false);
  });

  it('should update value when text is entered', async () => {
    render(GrowingTextArea, {});

    const textarea = page.getByRole('textbox');
    await userEvent.fill(textarea, 'Test content');

    await expect.element(textarea).toHaveValue('Test content');
  });

  it('should render with initial value', async () => {
    render(GrowingTextArea, {
      value: 'Initial text',
    });

    const textarea = page.getByRole('textbox');
    await expect.element(textarea).toHaveValue('Initial text');
  });

  it('should handle multiline text', async () => {
    render(GrowingTextArea, {});

    const textarea = page.getByRole('textbox');
    const multilineText = 'Line 1\nLine 2\nLine 3';
    await userEvent.fill(textarea, multilineText);

    await expect.element(textarea).toHaveValue(multilineText);
  });

  it('should call onkeydown handler when key is pressed', async () => {
    const onkeydown = vi.fn();
    render(GrowingTextArea, {
      onkeydown,
    });

    const textarea = page.getByRole('textbox');
    await userEvent.click(textarea);
    await userEvent.keyboard('a');

    expect(onkeydown).toHaveBeenCalled();
  });

  it('should handle Enter key press', async () => {
    const onkeydown = vi.fn();
    render(GrowingTextArea, {
      onkeydown,
    });

    const textarea = page.getByRole('textbox');
    await userEvent.click(textarea);
    await userEvent.keyboard('{Enter}');

    expect(onkeydown).toHaveBeenCalled();
  });

  it('should grow in height when content is added', async () => {
    const { container } = render(GrowingTextArea, {});

    const textarea = container.querySelector('textarea') as HTMLTextAreaElement;

    await userEvent.fill(page.getByRole('textbox'), 'Line 1\nLine 2\nLine 3\nLine 4\nLine 5');

    // Wait for height adjustment
    await new Promise(resolve => setTimeout(resolve, 100));

    // Verify the style.height has been set by the component
    expect(textarea.style.height).not.toBe('');
  });

  it('should handle empty value', async () => {
    render(GrowingTextArea, {
      value: '',
    });

    const textarea = page.getByRole('textbox');
    await expect.element(textarea).toHaveValue('');
  });

  it('should update when value changes multiple times', async () => {
    render(GrowingTextArea, {});

    const textarea = page.getByRole('textbox');

    await userEvent.fill(textarea, 'First text');
    await expect.element(textarea).toHaveValue('First text');

    await userEvent.clear(textarea);
    await userEvent.fill(textarea, 'Second text');
    await expect.element(textarea).toHaveValue('Second text');
  });

  it('should have rows attribute set to 1', () => {
    const { container } = render(GrowingTextArea, {});

    const textarea = container.querySelector('textarea');
    expect(textarea?.getAttribute('rows')).toBe('1');
  });

  it('should render with custom class', () => {
    const { container } = render(GrowingTextArea, {
      class: 'custom-class',
    });

    const textarea = container.querySelector('textarea');
    expect(textarea?.classList.contains('custom-class')).toBe(true);
  });

  it('should handle rapid text input', async () => {
    render(GrowingTextArea, {});

    const textarea = page.getByRole('textbox');
    await userEvent.fill(textarea, 'Quick');
    await userEvent.fill(textarea, 'Quick brown');
    await userEvent.fill(textarea, 'Quick brown fox');

    await expect.element(textarea).toHaveValue('Quick brown fox');
  });

  it('should be focusable', async () => {
    render(GrowingTextArea, {});

    const textarea = page.getByRole('textbox');
    await userEvent.click(textarea);

    const element = textarea.element() as HTMLTextAreaElement;
    expect(document.activeElement).toBe(element);
  });

  it('should focus when focus method is called', async () => {
    const { component } = render(GrowingTextArea, {});

    component.focus();

    const textarea = page.getByRole('textbox');
    const element = textarea.element() as HTMLTextAreaElement;
    expect(document.activeElement).toBe(element);
  });

  it('should reset height when resetHeight method is called', async () => {
    const { component, container } = render(GrowingTextArea, {});

    const textarea = container.querySelector('textarea') as HTMLTextAreaElement;

    // Add content to grow the textarea
    await userEvent.fill(page.getByRole('textbox'), 'Line 1\nLine 2\nLine 3\nLine 4\nLine 5');
    await new Promise(resolve => setTimeout(resolve, 100));

    // Get the grown height
    const grownHeight = textarea.style.height;
    expect(grownHeight).not.toBe('');

    // Call resetHeight
    component.resetHeight();

    // Height should be reset
    expect(textarea.style.height).toBe('');
  });

  it('should recalculate height after value changes', async () => {
    const { container } = render(GrowingTextArea, {});

    const textarea = container.querySelector('textarea') as HTMLTextAreaElement;

    // Add multiline content
    await userEvent.fill(page.getByRole('textbox'), 'Line 1\nLine 2\nLine 3');
    await new Promise(resolve => setTimeout(resolve, 100));

    // Verify height style was set
    expect(textarea.style.height).not.toBe('');
  });

  it('should handle input event when initialHeight is already set', async () => {
    const { component, container } = render(GrowingTextArea, {});

    const textarea = page.getByRole('textbox');
    const textareaElement = container.querySelector('textarea') as HTMLTextAreaElement;

    // Reset to ensure we start with initialHeight === 0
    component.resetHeight();

    // First input to set initialHeight (covers the if branch where initialHeight === 0)
    await userEvent.fill(textarea, 'First');
    await new Promise(resolve => setTimeout(resolve, 50));

    expect(textareaElement.style.height).not.toBe('');

    // Second input when initialHeight is already set (covers the else/skip branch where initialHeight !== 0)
    await userEvent.clear(textarea);
    await userEvent.fill(textarea, 'First\nSecond\nThird');
    await new Promise(resolve => setTimeout(resolve, 50));

    // Verify height is still being managed
    expect(textareaElement.style.height).not.toBe('');
  });
});
