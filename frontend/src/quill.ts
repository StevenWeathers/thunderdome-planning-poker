import Quill from 'quill'

function quill(node, options) {
    const quill = new Quill(node, {
        modules: {
            toolbar: [
                ['bold', 'italic', 'underline', 'strike'], // toggled buttons
                ['blockquote', 'code-block'],
                [{ list: 'ordered' }, { list: 'bullet' }],
                [{ script: 'sub' }, { script: 'super' }], // superscript/subscript
                [{ indent: '-1' }, { indent: '+1' }], // outdent/indent
                [{ header: [1, 2, 3, 4, 5, 6, false] }],
                [{ color: [] }, { background: [] }], // dropdown with defaults from theme
                [{ align: [] }],
                ['clean'], // remove formatting button
            ],
            clipboard: {
                matchVisual: false,
            },
        },
        placeholder: 'Type something...',
        theme: 'snow',
        ...options,
    })
    const container = node.getElementsByClassName('ql-editor')[0]
    if (options.content !== '') {
        const delta = quill.clipboard.convert(options.content)
        quill.setContents(delta)
    }

    quill.on('text-change', function (delta, oldDelta, source) {
        node.dispatchEvent(
            new CustomEvent('text-change', {
                detail: {
                    html: container.innerHTML.replace('<p><br></p>', ''),
                    text: quill.getText(),
                },
            }),
        )
    })
}

export { quill }
