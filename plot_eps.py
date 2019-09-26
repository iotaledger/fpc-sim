import numpy as np
import math
import sys
import csv
import matplotlib.pyplot as plt
from matplotlib import animation
import matplotlib

from mpl_toolkits.mplot3d import axes3d, Axes3D  # <-- Note the capitalization!
from matplotlib import pyplot as plt
folder = 'data/'


def main():
    # prepare data
    Z = loaddata()
    Z = Z[:, 1:Z.shape[1]]
    xlen = Z.shape[0]
    ylen = Z.shape[1]
    x = np.linspace(0, xlen, xlen)
    y = np.linspace(0, ylen, ylen)/ylen
    X, Y = np.meshgrid(x, y)
    Z = np.rot90(Z)
    max = 10000
    Z[Z > max] = max

    # 2D figure
    fig = plt.figure()
    ax = plt.axes()
    cmap = 'afmhot'
    cmap = matplotlib.cm.get_cmap(cmap+'_r')
    # cp = plt.contourf(X, Y, Z)
    cp = ax.imshow(Z, extent=[x[0], x[-1],  y[0], y[-1]],
                   cmap=cmap, aspect='auto')

    plt.xlabel("Time")
    plt.ylabel("$\\eta$")
    plt.title("Honest nodes per eta value")
    plt.colorbar(cp)
    # cbar = fig.colorbar(cax, ticks=[0.8, 0.9, 1])
    # cbar = fig.colorbar(cax)
    # cbar.ax.set_yticklabels(['<0.8', '0.9', '1.0'])
    plt.savefig(folder+'epsHisto.eps', format='eps')
    plt.clf()

    # # 3D figure
    # fig = plt.figure()
    # ax = fig.gca(projection='3d')
    # # ax.plot_surface(X, Y, Z, rstride=8, cstride=8, alpha=0.3)
    # # mycmap = plt.get_cmap('gist_earth')
    # # ax.plot_surface(X, Y, Z, cmap=mycmap)
    # # ax.contour3D(X, Y, Z, 50, cmap='binary')
    # ax.plot_surface(X, Y, Z, rstride=1, cstride=1,
    #                 cmap='gnuplot', edgecolor='none')
    # ax.view_init(30, 20)
    # plt.savefig(folder+'epsHisto3D.eps', format='eps')
    # plt.show()

    # # animated figure
    # fig = plt.figure()
    # ax = Axes3D(fig)

    # def animate(i):
    #     ax.view_init(30, i/2.-30)
    #     return fig,

    # def init():
    #     ax.plot_surface(X, Y, Z, rstride=1, cstride=1,
    #                     cmap='viridis', edgecolor='none')
    #     # ax.scatter(X, Y, Z, marker='o', s=20, c="goldenrod", alpha=0.6)
    #     return fig,
    # # Animate
    # anim = animation.FuncAnimation(
    #     fig, animate, init_func=init, frames=120, interval=20, blit=True)
    # # Save
    # anim.save('basic_animation.html', fps=30,
    #           extra_args=['-vcodec', 'libx264'])
    # plt.clf()


def loaddata():
    try:
        filestr = folder+'result_eta.csv'
        f = open(filestr, "r")
        data = np.loadtxt(filestr, delimiter=",")
        # Z = Z[1:Z.shape[0], :]
        # print(data)

        return data
    except FileNotFoundError:
        print(filestr)
        print("File not found.")
        return []


# define this at the very end of the file, it liberates the order of writing functions a bit
if __name__ == '__main__':
    main()
