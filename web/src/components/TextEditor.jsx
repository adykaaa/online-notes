import SunEditor from "suneditor-react";
import plugins from "suneditor/src/plugins";
import { en } from "suneditor/src/lang";
import katex from "katex";
import "suneditor/dist/css/suneditor.min.css";
import "katex/dist/katex.min.css";

const TextEditor = ({ name, onChange, props }) => {
  const options = {
    plugins: plugins,
    minHeight: "100px",
    maxHeight: "200px",
    minWidth: "200px",
    maxWidth: "800px",
    katex: katex,
    lang: en,
    buttonList: [
      [
        "font",
        "fontSize",
        "formatBlock",
        "bold",
        "underline",
        "italic",
        "paragraphStyle",
        "blockquote",
        "strike",
        "subscript",
        "superscript",
        "fontColor",
        "hiliteColor",
        "textStyle",
        "removeFormat",
        "undo",
        "redo",
        "outdent",
        "indent",
        "align",
        "horizontalRule",
        "list",
        "lineHeight",
        "table",
        "link",
        "image",
        'video',
        // 'audio',
        // You must add the 'katex' library at options to use the 'math' plugin.
        // 'math',
        // You must add the "imageGalleryUrl".
        // 'imageGallery',
        "fullScreen",
        "showBlocks",
        "codeView",
        "preview"
        // 'print'
        // 'save',
        // 'template'
      ]
    ]
  };


  return (
    <SunEditor
      {...props}
      placeholder="Please type here..."
      name={name}
      lang="en"
      setDefaultStyle="font-family: Arial; font-size: 14px;"
      setOptions={options}
      // onImageUpload={onImageUpload}
      onChange={onChange}
    />
  );
};

export default TextEditor;
