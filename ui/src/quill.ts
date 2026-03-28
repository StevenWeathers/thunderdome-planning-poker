import Quill from 'quill';

function quill(node, options) {
  let quill = createQuill(node, options);

  return {
    update(newOptions) {
      // handle reactive updates if needed
      // e.g. placeholder/content changes
      if (newOptions.content !== options.content) {
        const delta = quill.clipboard.convert(newOptions.content);
        quill.setContents(delta);
      }
      options = newOptions;
    },

    destroy() {
      quill.off('text-change');
      // optionally clear DOM
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
