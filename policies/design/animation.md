# Animation & Motion Policy

## Forbidden Excessive Animation Slop

Agents MUST NOT automatically add:
- Animated gradient backgrounds
- Floating card bounce effects
- Parallax background scrolling everywhere
- Text reveal animations on every header
- Hover scale animations on every element
- Continuous background noise animations

## Purposeful Motion Rules

1. **State Feedback Only**: Animation must serve a functional purpose (modal opening, dropdown transition, form error feedback).
2. **Reduced Motion Support**: Always respect `prefers-reduced-motion` media queries.
3. **Subtle Durations**: Keep transition durations subtle (150ms - 250ms).
