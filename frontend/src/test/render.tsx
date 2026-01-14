import { type ReactElement } from "react";
import { MemoryRouter, useLocation } from "react-router-dom";
import { render } from "@testing-library/react";

export function LocationDisplay() {
  const loc = useLocation();
  return <div data-testid="location">{loc.pathname + loc.search}</div>;
}

export function renderWithRoute(ui: ReactElement, route: string) {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <LocationDisplay />
      {ui}
    </MemoryRouter>
  );
}