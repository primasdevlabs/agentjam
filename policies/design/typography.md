# Typography Policy

## Typography Rules

1. **Do Not Default to AI Typography Slop**: Do NOT automatically select Inter, Roboto, Poppins, Montserrat, Space Grotesk, or browser fallback defaults simply because no font is specified.
2. **Inspect Existing System**:
   - Inspect existing project typography declarations, CSS `@import` / Google Font loads, and CSS custom properties.
   - Preserve existing heading/body hierarchy.
3. **Request Direction for New Projects**:
   - For new projects with no typography system, request typography direction before picking fonts.
4. **Explicit "Noto" Exception**:
   - If the user explicitly specifies **"use Noto"**, use the appropriate Noto font family (e.g., Noto Sans / Noto Serif) and do not request typography clarification again.
