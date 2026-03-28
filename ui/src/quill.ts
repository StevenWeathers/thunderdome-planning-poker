import Quill from 'quill';

function quill(node, options) {
  let quill = createQuill(node, options);
  let isLocalChange = false;

  // Patch text-change to set isLocalChange
  quill.on('text-change', () => {
    isLocalChange = true;
  });

  return {
    update(newOptions) {
      // Only update if content changed externally
      const currentHtml = node.getElementsByClassName('ql-editor')[0]?.innerHTML || '';
      // Remove <p><br></p> for comparison, as in event
      const normalizedCurrent = currentHtml.replace('<p><br></p>', '');
      const normalizedNew = (newOptions.content || '').replace('<p><br></p>', '');
      if (!isLocalChange && normalizedNew !== normalizedCurrent) {
        const delta = quill.clipboard.convert(newOptions.content || '');
        quill.setContents(delta);
      }
      isLocalChange = false;
      options = newOptions;
    },

    destroy() {
      quill.off('text-change');
      node.innerHTML = '';
    },
  };
}

function createQuill(node, options) {
  const quill = new Quill(node, {
    modules: {
      /* ... */
    },
    theme: 'snow',
    ...options,
  });

  const container = node.getElementsByClassName('ql-editor')[0];

  if (options.content) {
    const delta = quill.clipboard.convert(options.content);
    quill.setContents(delta);
  }

  quill.on('text-change', () => {
    node.dispatchEvent(
      new CustomEvent('textchange', {
        detail: {
          html: container.innerHTML.replace('<p><br></p>', ''),
          text: quill.getText(),
        },
      }),
    );
  });

  return quill;
}

export { quill };
