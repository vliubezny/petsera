import { mount } from "@cypress/vue";
import { reactive, ref } from "vue";
import PetMap from "@/components/PetMap.vue";

const createPetMap = (opts = {}) =>
  mount(PetMap, {
    propsData: {
      ...opts,
    },
    styles: `
  body{
    margin: 0;
  }
  [data-testid=map]{
    height: 100vh;
    width: 100vw;
  }
  `,
  });

it("should load map", () => {
  createPetMap();

  cy.get("[aria-label='Show street map']").should("be.visible");
});

it("should show markers", () => {
  const markers = reactive([
    {
      position: {
        lat: 53.908,
        lng: 27.528,
      },
      title: "Marker1",
    },
    {
      position: {
        lat: 53.886,
        lng: 27.563,
      },
      title: "Marker2",
    },
  ]);
  createPetMap({
    markers,
  });

  cy.get("div[role=button]").should("have.length", 2);

  cy.then(() => {
    markers.push({
      position: {
        lat: 53.896,
        lng: 27.544,
      },
      title: "Marker3",
    });
  });

  cy.get("div[role=button]").should("have.length", 3);
});

it("should emit bounds-changed events", () => {
  createPetMap();

  expect(Cypress.vueWrapper.emitted()["bounds-changed"]).has.lengthOf(1);
  cy.get('[aria-label="Zoom in"]').click();

  cy.then(() => {
    expect(Cypress.vueWrapper.emitted()["bounds-changed"]).has.lengthOf(2);
  });
});

it("should emit marker-selected events", () => {
  const markers = reactive([
    {
      position: {
        lat: 53.908,
        lng: 27.528,
      },
      title: "Marker1",
    },
  ]);
  createPetMap({
    markers,
  });

  cy.get("div[role=button]").click();

  cy.then(() => {
    const events = Cypress.vueWrapper.emitted()["marker-selected"];
    expect(events).has.lengthOf(1);
    expect(events[0][0]).to.be.deep.equal(markers[0]);
  });
});

it.only("should select position", () => {
  const selectable = ref(false);
  createPetMap({
    selectable,
  });

  cy.get("[aria-label='Map']").click();
  cy.then(() => {
    expect(Cypress.vueWrapper.emitted()["position-selected"]).to.not.be.ok;
  });

  cy.then(() => {
    selectable.value = true;
  });

  cy.wr;

  cy.then(() => Cypress.vue.$nextTick().then(() => console.log("tick"))).then(
    () => console.log("then")
  );

  cy.get("[aria-label='Map']").click();

  cy.get("[title='New position']").should("be.visible");
  cy.then(() => {
    expect(Cypress.vueWrapper.emitted()["position-selected"]).has.lengthOf(1);
  });
});
