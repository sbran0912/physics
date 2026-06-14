// main.ts – analog zu main.go (Port April 2026)

import * as std from "./lib-2d.ts";
import * as phy from "./lib-phys.ts";

const SCREEN_W = 1400;
const SCREEN_H = 700;

let shapes: phy.Shape[] = [];

function start() {
  // Boden (unbeweglich)
  shapes.push(new phy.Polygon(350, 500, 1000, 50, true));

  // Mauer aus 4 gestapelten Blöcken mit zufälliger Breite und Höhe
  const baseX = 800;
  const baseY = 500;
  let currentY = baseY;

  for (let i = 0; i < 4; i++) {
    const blockW = 60 + Math.random() * 40; // Breite zwischen 60 und 100
    const blockH = 35 + Math.random() * 35; // Höhe zwischen 35 und 70
    const blockX = baseX - blockW / 2; // Zentriert

    shapes.push(
      new phy.Polygon(blockX, currentY - blockH, blockW, blockH, false),
    );
    currentY -= blockH; // Für nächsten Block
  }

  // Zwei Kreise
  const circle1 = new phy.Circle(400, 400, 40, false);
  const circle2 = new phy.Circle(300, 300, 35, false);
  shapes.push(circle1, circle2);

  // Einmaliger Impuls nach rechts
  circle1.applyForce(new std.Vector(15, 0), 0);
  circle2.applyForce(new std.Vector(20, 0), 0);
}

function draw() {
  std.background("#0B1A2E"); // DarkBlue-Ersatz

  // Gravität + State-Reset
  for (const s of shapes) {
    s.applyGravity(new std.Vector(0, 0.4));
    s.clearState();
  }

  // Kollisionen erkennen und auflösen
  for (let i = 0; i < shapes.length - 1; i++) {
    for (let j = i + 1; j < shapes.length; j++) {
      const { isColliding, mtv, contacts } = phy.detectCollision(
        shapes[i],
        shapes[j],
      );
      if (isColliding) {
        const { forceA, angForceA, forceB, angForceB } = phy.resolveCollision(
          shapes[i],
          shapes[j],
          contacts,
          mtv,
        );
        shapes[i].applyForce(forceA, angForceA);
        shapes[j].applyForce(forceB, angForceB);
      }
    }
  }

  // Update & Draw
  const color = "#F5F0DC"; // Beige-Ersatz
  for (const s of shapes) {
    s.update(0.8);
    s.draw(color, 2);
  }
}

std.init(SCREEN_W, SCREEN_H);
start();
std.startAnimation(draw);
