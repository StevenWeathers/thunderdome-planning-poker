/// <reference types="svelte" />
/// <reference types="vite/client" />

declare module 'snapsvg-cjs' {
  interface Element {
    select(query: string): Element;
    attr(name: string): any;
    attr(name: string, value: any): Element;
    attr(params: object): Element;
  }

  interface Animation {
    stop(): void;
  }

  interface Snap {
    (element: HTMLElement | SVGElement | string): Element;
    select(query: string): Element;
    animate(
      from: number,
      to: number,
      setter: (val: number) => void,
      duration: number,
      easing?: (n: number) => number,
      callback?: () => void,
    ): Animation;
  }

  interface Mina {
    easeinout: (n: number) => number;
    linear: (n: number) => number;
    easein: (n: number) => number;
    easeout: (n: number) => number;
  }

  const Snap: Snap;
  export default Snap;
  export const mina: Mina;
}
