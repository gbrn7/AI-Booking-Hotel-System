import React from "react";

type FormatType = "code" | "bold" | "underline" | "strike" | "italic";

interface FormatRule {
  type: FormatType;
  regex: RegExp;
}

const RULES: FormatRule[] = [
  { type: "code", regex: /`([^`\n]+)`/ },
  { type: "bold", regex: /\*\*([^\s*](?:[\s\S]*?[^\s*])?)\*\*/ },
  { type: "bold", regex: /<(?:b|strong)>([\s\S]*?)<\/(?:b|strong)>/i },
  { type: "underline", regex: /__([^\s_](?:[\s\S]*?[^\s_])?)__/ },
  { type: "underline", regex: /<(?:u|ins)>([\s\S]*?)<\/(?:u|ins)>/i },
  { type: "strike", regex: /~~([^\s~](?:[\s\S]*?[^\s~])?)~~/ },
  { type: "strike", regex: /(?<!~)\~([^\s~](?:[\s\S]*?[^\s~])?)\~(?!~)/ },
  { type: "strike", regex: /<(?:s|del|strike)>([\s\S]*?)<\/(?:s|del|strike)>/i },
  { type: "italic", regex: /(?<!\*)\*([^\s*](?:[\s\S]*?[^\s*])?)\*(?!\*)/ },
  { type: "italic", regex: /(?<!_)_([^\s_](?:[\s\S]*?[^\s_])?)_(?!_)/ },
  { type: "italic", regex: /<(?:i|em)>([\s\S]*?)<\/(?:i|em)>/i },
];

/**
 * Formats chat message text by converting indicators for:
 * - Bold: **text**, <b>text</b>, <strong>text</strong>
 * - Underline: __text__, <u>text</u>, <ins>text</ins>
 * - Strikethrough (teks coret): ~~text~~, ~text~, <s>text</s>, <del>text</del>, <strike>text</strike>
 * - Italic: *text*, _text_, <i>text</i>, <em>text</em>
 * - Code: `text`
 *
 * Supports nested formatting safely without dangerous HTML injection.
 */
export function formatChatText(input: string, keyPrefix = "fmt", depth = 0): React.ReactNode[] {
  if (!input || depth > 10) return [input];

  let bestMatch: RegExpExecArray | null = null;
  let bestRule: FormatRule | null = null;

  for (const rule of RULES) {
    const match = rule.regex.exec(input);
    if (match) {
      if (!bestMatch || match.index < bestMatch.index) {
        bestMatch = match;
        bestRule = rule;
      }
    }
  }

  if (!bestMatch || !bestRule) {
    return [input];
  }

  const before = input.slice(0, bestMatch.index);
  const matchedText = bestMatch[0];
  const innerContent = bestMatch[1];
  const after = input.slice(bestMatch.index + matchedText.length);

  const result: React.ReactNode[] = [];

  if (before) {
    result.push(before);
  }

  const childKey = `${keyPrefix}-${bestMatch.index}`;
  const children = formatChatText(innerContent, `${childKey}-in`, depth + 1);

  switch (bestRule.type) {
    case "bold":
      result.push(
        <strong key={childKey} className="font-semibold text-inherit">
          {children}
        </strong>
      );
      break;
    case "underline":
      result.push(
        <u key={childKey} className="underline underline-offset-2 decoration-current">
          {children}
        </u>
      );
      break;
    case "strike":
      result.push(
        <s key={childKey} className="line-through decoration-current opacity-75">
          {children}
        </s>
      );
      break;
    case "italic":
      result.push(
        <em key={childKey} className="italic">
          {children}
        </em>
      );
      break;
    case "code":
      result.push(
        <code key={childKey} className="px-1 py-0.5 rounded text-[0.875em] font-mono bg-current/10">
          {children}
        </code>
      );
      break;
    default:
      result.push(children);
  }

  if (after) {
    result.push(...formatChatText(after, `${keyPrefix}-${bestMatch.index}-post`, depth));
  }

  return result;
}
