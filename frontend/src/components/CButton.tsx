import { Button } from "@/components/ui/button";

export const CButton = ({ ...props }) => {
  return <Button className="cursor-pointer" size="lg" {...props}></Button>;
};
