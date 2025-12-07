import { useEffect, useRef, useState } from "react";
import { Editor } from "@tiptap/react";
import { Bold, Heading1, Italic } from "lucide-react";
import { CommandType } from "../types";

const CommandMenu = ({
  editor,
  keyDown,
}: {
  editor: Editor;
  keyDown: string;
}) => {
  const [isMenuOpen, setMenuOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [filteredCommands, setFilteredCommands] = useState<CommandType[]>([]);
  const [menuPosition, setMenuPosition] = useState({ top: 0, left: 0 });

  const commands = [
    {
      name: "Heading",
      icon: <Heading1 size={18} />,
      action: () => editor.chain().focus().toggleHeading({ level: 1 }).run(),
    },
    {
      name: "Bold",
      icon: <Bold size={18} />,
      action: () => editor.chain().focus().toggleBold().run(),
    },
    {
      name: "Italic",
      icon: <Italic size={18} />,
      action: () => editor.chain().focus().toggleItalic().run(),
    },

    // Add more commands as needed
  ];

  const menuRef = useRef(null);

  useEffect(() => {
    if (keyDown === "/") {
      setMenuOpen(true);
      setQuery("/");
    } else if (isMenuOpen) {
      if (keyDown === "Escape") {
        setMenuOpen(false);
        setQuery("");
      } else if (keyDown.length === 1) {
        setQuery((prev) => prev + keyDown);
      }
    }
  }, [keyDown]);

  useEffect(() => {
    if (query.startsWith("/")) {
      const search = query.slice(1).toLowerCase();
      const matches = commands.filter((command) =>
        command.name.toLowerCase().includes(search)
      );
      setFilteredCommands(matches);

      // Calculate menu position
      const { from } = editor.state.selection;
      const coords = editor.view.coordsAtPos(from);

      setMenuPosition({
        top: coords.top + window.scrollY + 25,
        left: coords.left + window.scrollX,
      });
    } else {
      setFilteredCommands([]);
    }
  }, [query]);

  const handleCommandSelect = (command: CommandType) => {
    command.action();
    setMenuOpen(false);
    setQuery("");
  };

  return (
    <div>
      {isMenuOpen && (
        <div
          ref={menuRef}
          className="shadow-md min-h-8 rounded-md w-44 absolute bg-white z-10 border"
          style={{
            position: "absolute",
            top: `${menuPosition.top}px`,
            left: `${menuPosition.left}px`,
          }}
        >
          {filteredCommands.map((command, index) => (
            <div
              key={index}
              className={`command-item hover:bg-neutral-100 px-4 py-1 cursor-pointer flex items-center gap-x-4 border-b border-neutral-200`}
              onClick={() => handleCommandSelect(command)}
            >
              {command.icon}
              {command.name}
            </div>
          ))}

          {filteredCommands.length === 0 && (
            <div className="px-4 py-1 text-neutral-500">No commands found</div>
          )}
        </div>
      )}
    </div>
  );
};

export default CommandMenu;
