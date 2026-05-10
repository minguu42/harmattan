import { Dialog as DialogBase } from "@base-ui/react";
import type { Dispatch, ReactNode, SetStateAction } from "react";

type Props = {
  children: ReactNode;
  open: boolean;
  setOpen: Dispatch<SetStateAction<boolean>>;
};

export function Dialog({ children, open, setOpen }: Props) {
  return (
    <DialogBase.Root open={open} onOpenChange={setOpen}>
      <DialogBase.Portal>
        <DialogBase.Backdrop className="fixed inset-0 bg-scrim" />
        <DialogBase.Popup className="fixed top-1/2 left-1/2 -mt-32 flex w-400 -translate-x-1/2 -translate-y-1/2 flex-col gap-16 rounded-xl border border-border bg-background p-16">
          {children}
        </DialogBase.Popup>
      </DialogBase.Portal>
    </DialogBase.Root>
  );
}

type DialogTitleProps = {
  children: string;
};

export function DialogTitle({ children }: DialogTitleProps) {
  return (
    <DialogBase.Title className="text-base font-medium text-on-background">
      {children}
    </DialogBase.Title>
  );
}
