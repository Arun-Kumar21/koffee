import { useCallback, useState } from "react";
import * as Y from "yjs";
import { Editor as TipTapEditor } from "@tiptap/core";
import { useEditor } from "@tiptap/react";
import Collaboration from "@tiptap/extension-collaboration";
import StarterKit from "@tiptap/starter-kit";
import { Socket } from "socket.io-client";
import { decodeUpdate, encodeUpdate } from "../../../lib/editor-helper";

export const useCollaborativeEditor = (
  socket: Socket | null,
) => {
  const [ydoc] = useState(() => new Y.Doc());
  const [isInitialized, setIsInitialized] = useState(false);

  const handleEditorCreate = useCallback(
    ({ editor }: { editor: TipTapEditor }) => {
      if (!isInitialized) {
        editor.commands.setContent("");
        setIsInitialized(true);
      }
    },
    [isInitialized]
  );

  const handleDocUpdate = useCallback(() => {
    if (!socket) return;
    const update = Y.encodeStateAsUpdate(ydoc);
    const encodedUpdate = encodeUpdate(update);
    socket.emit("update-canvas", encodedUpdate);
  }, [socket, ydoc]);

  const handleUpdate = useCallback(
    (update: string) => {
      try {
        const updateArray = decodeUpdate(update);
        Y.applyUpdate(ydoc, updateArray);
      } catch (error) {
        console.error("Error applying update:", error);
      }
    },
    [ydoc]
  );

  const handleInitialState = useCallback(
    (initialState: string) => {
      try {
        const stateArray = decodeUpdate(initialState);
        Y.applyUpdate(ydoc, stateArray);
      } catch (error) {
        console.error("Error applying initial state:", error);
      }
    },
    [ydoc]
  );

  const updateHandler = useCallback(
    (update: Uint8Array) => {
      if (!socket || update.byteLength === 0) return;
      const encodedUpdate = encodeUpdate(update);
      socket.emit("update-canvas", encodedUpdate);
    },
    [socket]
  );

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        history: false,
      }),
      Collaboration.configure({
        document: ydoc,
        field: "content",
      }),
    ],
    onCreate: handleEditorCreate,
    onUpdate: handleDocUpdate,
    editorProps: {
      attributes: {
        class:
          "h-full mx-auto focus:outline-none prose prose-sm sm:prose lg:prose-lg whitespace-pre-wrap break-words overflow-auto",
      },
    },
  });

  return {
    editor,
    ydoc,
    handleUpdate,
    handleInitialState,
    updateHandler,
  };
};
